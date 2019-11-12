package service

import (
	"fmt"

	"gitlab.com/diancai/diancai-api/database"
	"gitlab.com/diancai/diancai-api/model"
	common "gitlab.com/diancai/diancai-services-common"
	"gitlab.com/diancai/diancai-services-common/auth"
	"gitlab.com/diancai/diancai-services-common/derror"
)

type ShopService interface {
	GetShop(doer interface{}, id string) (model.Shop, error)
	GetCategoriesByShop(doer interface{}, shopID string) ([]model.Category, error)
	GetMenusByCategory(doer interface{}, shopID string, categoryID string) ([]model.MenuItem, error)
	CreateShop(doer interface{}, shop model.Shop) (model.Shop, error)
	CreateMenu(doer interface{}, shopID string, menu model.MenuItem) (model.MenuItem, error)
	EditShop(doer interface{}, id string, shop model.Shop) (model.Shop, error)
	EditShopMenu(doer interface{}, shopID string, menuID string, menu model.MenuItem) (model.MenuItem, error)
	EditShopCategories(doer interface{}, id string, categories []model.Category) ([]model.Category, error)
	EditShopCategory(doer interface{}, shopID string, categoryID string, category model.Category) (model.Category, error)
	DeleteShop(doer interface{}, id string) error
	DeleteMenu(doer interface{}, shopID string, menuID string) error
}

type shopServiceImpl struct {
	db   database.ShopDatabase
	time common.TimeSource
}

func ShopServiceInit(db database.ShopDatabase) ShopService {
	return &shopServiceImpl{
		db:   db,
		time: &common.RealTime{},
	}
}

func (s *shopServiceImpl) GetShop(doer interface{}, id string) (model.Shop, error) {
	if id == "" {
		return model.Shop{}, derror.ErrorDebug(common.ErrInputValidationFailed, "get shop: blank ID")
	}
	shop, err := s.db.GetShop(id)
	if err != nil {
		return model.Shop{}, derror.ErrorDebug(err, fmt.Sprintf("get shop: shop=%v", shop))
	}
	return shop, nil
}

func (s *shopServiceImpl) GetCategoriesByShop(doer interface{}, shopID string) ([]model.Category, error) {
	if shopID == "" {
		return nil, derror.ErrorDebug(common.ErrInputValidationFailed, "get categories by shop: blank ID")
	}
	categories, err := s.db.GetCategoriesByShop(shopID)
	if err != nil {
		return nil, derror.ErrorDebug(err, fmt.Sprintf("get categories by shop: shopID=%s", shopID))
	}
	return categories, nil
}

func (s *shopServiceImpl) GetMenusByCategory(doer interface{}, shopID string, categoryID string) ([]model.MenuItem, error) {
	if shopID == "" || categoryID == "" {
		return nil, derror.Error(common.ErrInputValidationFailed)
	}
	cat, err := s.db.GetCategory(categoryID)
	if err != nil {
		return nil, derror.Error(err)
	}
	menus, err := s.db.GetMenus(cat.Menu)
	if err != nil {
		return nil, derror.Error(err)
	}
	for i := range menus {
		if err := menus[i].Sign(); err != nil {
			return nil, derror.ErrorDebug(err, "get menu: error signing menu")
		}
	}
	return menus, nil
}

func (s *shopServiceImpl) CreateShop(doer interface{}, shop model.Shop) (model.Shop, error) {
	if shop.ID != "" {
		return model.Shop{}, derror.Error(common.ErrInputValidationFailed)
	}

	if !model.ACLCheckRoleIsAdmin(doer) {
		return model.Shop{}, derror.ErrorDebug(common.ErrForbidden, fmt.Sprintf("create shop failed: forbidden: doer is not admin: %v", doer))
	}

	currentTime := s.time.Now()
	shop.CreatedAt = currentTime
	shop.UpdatedAt = currentTime

	shop, err := s.db.CreateShop(shop)
	if err != nil {
		return model.Shop{}, derror.Error(err)
	}

	// add default category
	_, err = s.db.CreateCategoriesForShop(shop.ID, []model.Category{model.DefaultCategory})
	if err != nil {
		return model.Shop{}, derror.Error(err)
	}

	return shop, nil
}

func (s *shopServiceImpl) CreateMenu(doer interface{}, shopID string, menu model.MenuItem) (model.MenuItem, error) {
	if shopID == "" || menu.ID != "" {
		return model.MenuItem{}, derror.Error(common.ErrInputValidationFailed)
	}

	if !model.ACLCheckRoleByShopIdAndLevel(doer, shopID, auth.RoleShopOwner) {
		return model.MenuItem{}, derror.ErrorDebug(common.ErrForbidden, fmt.Sprintf("create menu failed: forbidden: doer is not authorised to create menu for shop id %s: %v", shopID, doer))
	}

	menu.ShopID = shopID
	if menu.OptionGroups == nil {
		menu.OptionGroups = make([]model.MenuOptionGroup, 0)
	}

	if menu.CategoryIDs == nil {
		menu.CategoryIDs = make([]string, 0)
	}

	curTime := s.time.Now()
	menu.CreatedAt = curTime
	menu.UpdatedAt = curTime

	menu, err := s.db.CreateMenu(menu)
	if err != nil {
		return model.MenuItem{}, derror.Error(err)
	}

	err = s.db.AppendMenuToCategories(shopID, menu.CategoryIDs, menu.ID)
	if err != nil {
		return model.MenuItem{}, derror.Error(err)
	}

	return menu, nil
}

func (s *shopServiceImpl) EditShop(doer interface{}, id string, shop model.Shop) (model.Shop, error) {
	if id != shop.ID || id == "" {
		return model.Shop{}, derror.Error(common.ErrInputValidationFailed)
	}

	if !model.ACLCheckRoleIsAdmin(doer) {
		return model.Shop{}, derror.ErrorDebug(common.ErrForbidden, fmt.Sprintf("edit shop failed: forbidden: doer is not admin: %v", doer))
	}

	shop.UpdatedAt = s.time.Now()
	shop, err := s.db.EditShop(id, shop)
	if err != nil {
		return model.Shop{}, derror.Error(err)
	}
	return shop, nil
}

func (s *shopServiceImpl) EditShopMenu(doer interface{}, shopID string, menuID string, menu model.MenuItem) (model.MenuItem, error) {
	if shopID == "" || menuID == "" {
		return model.MenuItem{}, derror.Error(common.ErrInputValidationFailed)
	}

	if !model.ACLCheckRoleByShopIdAndLevel(doer, shopID, auth.RoleShopOwner) {
		return model.MenuItem{}, derror.ErrorDebug(common.ErrForbidden, fmt.Sprintf("edit menu failed: forbidden: doer is not authorised to edit menu for shop id %s: %v", shopID, doer))
	}

	oldMenu, err := s.db.GetMenu(menuID)
	if err != nil {
		return model.MenuItem{}, derror.Error(err)
	}
	removeCat := map[string]bool{}
	addCat := map[string]bool{}
	removeCatTotal := 0
	addCatTotal := 0
	for _, cat := range oldMenu.CategoryIDs {
		removeCat[cat] = true
		removeCatTotal++
	}
	for _, cat := range menu.CategoryIDs {
		if !removeCat[cat] {
			addCat[cat] = true
			addCatTotal++
		}
		removeCat[cat] = false
		removeCatTotal--
	}
	menu.UpdatedAt = s.time.Now()
	menu, err = s.db.EditMenu(menuID, menu)
	if err != nil {
		return model.MenuItem{}, derror.Error(err)
	}
	removeCatList := make([]string, 0, removeCatTotal)
	addCatList := make([]string, 0, addCatTotal)
	for cat, add := range addCat {
		if !add {
			continue
		}
		addCatList = append(addCatList, cat)
	}
	for cat, rem := range removeCat {
		if !rem {
			continue
		}
		removeCatList = append(removeCatList, cat)
	}
	err = s.db.AppendMenuToCategories(shopID, addCatList, menuID)
	if err != nil {
		return model.MenuItem{}, derror.Error(err)
	}
	err = s.db.RemoveMenuFromCategories(shopID, removeCatList, menuID)
	if err != nil {
		return model.MenuItem{}, derror.Error(err)
	}
	return menu, nil
}

func (s *shopServiceImpl) EditShopCategories(doer interface{}, shopID string, categories []model.Category) ([]model.Category, error) {
	if shopID == "" || len(categories) == 0 {
		return nil, common.ErrInputValidationFailed
	}

	if !model.ACLCheckRoleByShopIdAndLevel(doer, shopID, auth.RoleShopOwner) {
		return nil, derror.ErrorDebug(common.ErrForbidden, fmt.Sprintf("edit categories failed: forbidden: doer is not authorised to edit categories for shop id %s: %v", shopID, doer))
	}

	// delete prev category
	err := s.db.DeleteCategoriesByShop(shopID)
	if err != nil {
		return nil, derror.Error(err)
	}

	for i := range categories {
		categories[i].Position = int64(i + 1)
	}

	// add category
	categories, err = s.db.CreateCategoriesForShop(shopID, categories)
	if err != nil {
		return nil, derror.Error(err)
	}

	return categories, nil
}

func (s *shopServiceImpl) EditShopCategory(doer interface{}, shopID string, categoryID string, category model.Category) (model.Category, error) {
	if shopID == "" || categoryID == "" {
		return model.Category{}, derror.Error(common.ErrInputValidationFailed)
	}

	if !model.ACLCheckRoleByShopIdAndLevel(doer, shopID, auth.RoleShopOwner) {
		return model.Category{}, derror.ErrorDebug(common.ErrForbidden, fmt.Sprintf("edit category failed: forbidden: doer is not authorised to edit category for shop id %s: %v", shopID, doer))
	}

	category, err := s.db.EditCategory(categoryID, category)
	if err != nil {
		return model.Category{}, derror.Error(err)
	}

	return category, nil
}

func (s *shopServiceImpl) DeleteShop(doer interface{}, id string) error {
	if id == "" {
		return derror.Error(common.ErrInputValidationFailed)
	}

	if !model.ACLCheckRoleIsAdmin(doer) {
		return derror.ErrorDebug(common.ErrForbidden, fmt.Sprintf("delete shop failed: forbidden: doer is not admin: %v", doer))
	}

	err := s.db.DeleteShop(id, s.time.Now())
	if err != nil {
		return derror.Error(err)
	}
	return nil
}

func (s *shopServiceImpl) DeleteMenu(doer interface{}, shopID string, menuID string) error {
	if shopID == "" || menuID == "" {
		return derror.Error(common.ErrInputValidationFailed)
	}

	if !model.ACLCheckRoleByShopIdAndLevel(doer, shopID, auth.RoleShopOwner) {
		return derror.ErrorDebug(common.ErrForbidden, fmt.Sprintf("delete menu failed: forbidden: doer is not authorised to delete menu for shop id %s: %v", shopID, doer))
	}

	err := s.db.RemoveMenuFromAllCategories(shopID, menuID)
	if err != nil {
		return derror.Error(err)
	}
	err = s.db.DeleteMenu(menuID, s.time.Now())
	if err != nil {
		return derror.Error(err)
	}
	return nil
}
