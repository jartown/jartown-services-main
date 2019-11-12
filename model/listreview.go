package model

import (
	"encoding/hex"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gitlab.com/diancai/diancai-api/config"
)

type Shop struct {
	ID           string     `bson:"-" json:"id"`
	Name         StringIntl `bson:"name" json:"name"`
	Description  StringIntl `bson:"description" json:"description"`
	Logo         string     `bson:"logo" json:"logo"`
	Photos       []string   `bson:"photos" json:"photos"`
	Phones       []string   `bson:"phones" json:"phones"`
	VerifyMethod string     `bson:"verify_method" json:"verify_method"`
	Address      StringIntl `bson:"address" json:"address"`
	Coordinates  []float64  `bson:"coordinates" json:"coordinates"`
	Tables       []string   `bson:"tables" json:"tables"`
	CreatedAt    time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt    time.Time  `bson:"deleted_at" json:"deleted_at"`
}

type MenuItem struct {
	ID           string            `bson:"-" json:"id"`
	Identifier   string            `bson:"identifier" json:"identifier"`
	ShopID       string            `bson:"shop_id" json:"shop_id"`
	CategoryIDs  []string          `bson:"category_ids" json:"category_ids"`
	Name         StringIntl        `bson:"name" json:"name"`
	Description  StringIntl        `bson:"description" json:"description"`
	Image        string            `bson:"image" json:"image"`
	Price        int64             `bson:"price" json:"price"`
	OptionGroups []MenuOptionGroup `bson:"option_groups" json:"option_groups"`
	CreatedAt    time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time         `bson:"updated_at" json:"updated_at"`
	DeletedAt    time.Time         `bson:"deleted_at" json:"deleted_at"`
	SigToken     string            `bson:"-" json:"sig_token,omitempty"`
}

type MenuOptionGroup struct {
	Name       StringIntl   `bson:"name" json:"name"`
	Options    []MenuOption `bson:"options" json:"options"`
	MinOptions int64        `bson:"min" json:"min"`
	MaxOptions int64        `bson:"max" json:"max"`
}

type MenuOption struct {
	Identifier string     `bson:"identifier" json:"identifier"`
	Name       StringIntl `bson:"name" json:"name"`
	ExtraPrice int64      `bson:"extra_price" json:"extra_price"`
}

type Category struct {
	ID       string     `bson:"-" json:"id"`
	ShopID   string     `bson:"shop_id" json:"shop_id"`
	Position int64      `bson:"position" json:"position"`
	Name     StringIntl `bson:"name" json:"name"`
	Menu     []string   `bson:"menu" json:"menu"`
}

var DefaultCategory = Category{
	Position: 1,
	Name: StringIntl{
		EN: "Menu",
		TH: "เมนู",
		ZH: "菜单",
	},
}

type MenuItemClaims struct {
	Name    StringIntl            `json:"name"`
	Price   int64                 `json:"price"`
	Options map[string]MenuOption `json:"options"`
	jwt.StandardClaims
}

func (m MenuItem) GetSignature() (string, error) {
	claims := MenuItemClaims{}
	claims.Name = m.Name
	claims.Price = m.Price
	claims.Options = make(map[string]MenuOption)
	for _, optGroup := range m.OptionGroups {
		for _, opt := range optGroup.Options {
			claims.Options[opt.Identifier] = opt
		}
	}
	claims.Subject = m.Identifier
	now := time.Now()
	expires := now.Add(time.Duration(config.Conf.JWTMenu.Timeout) * time.Second)
	claims.IssuedAt = now.UTC().Unix()
	claims.NotBefore = now.UTC().Unix()
	claims.ExpiresAt = expires.UTC().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	hexKey, err := hex.DecodeString(config.Conf.JWTMenu.Key)
	if err != nil {
		return "", err
	}
	jwtStr, err := token.SignedString(hexKey)
	if err != nil {
		return "", err
	}

	return jwtStr, nil
}

func (m *MenuItem) Sign() error {
	sig, err := m.GetSignature()
	if err != nil {
		return err
	}

	m.SigToken = sig

	return nil
}
