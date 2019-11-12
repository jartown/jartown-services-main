package model

import (
	"encoding/json"
)

func UnmarshalListReview(data []byte) (ListReview, error) {
	var r ListReview
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ListReview) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ListReview struct {
	ID                   string       `json:"_id"`
	ListingURL           string       `json:"listing_url"`
	Name                 string       `json:"name"`
	Summary              string       `json:"summary"`
	Space                string       `json:"space"`
	Description          string       `json:"description"`
	NeighborhoodOverview string       `json:"neighborhood_overview"`
	Notes                string       `json:"notes"`
	Transit              string       `json:"transit"`
	Access               string       `json:"access"`
	Interaction          string       `json:"interaction"`
	HouseRules           string       `json:"house_rules"`
	PropertyType         string       `json:"property_type"`
	RoomType             string       `json:"room_type"`
	BedType              string       `json:"bed_type"`
	MinimumNights        string       `json:"minimum_nights"`
	MaximumNights        string       `json:"maximum_nights"`
	CancellationPolicy   string       `json:"cancellation_policy"`
	LastScraped          string       `json:"last_scraped"`
	CalendarLastScraped  string       `json:"calendar_last_scraped"`
	FirstReview          string       `json:"first_review"`
	LastReview           string       `json:"last_review"`
	Accommodates         int64        `json:"accommodates"`
	Bedrooms             int64        `json:"bedrooms"`
	Beds                 int64        `json:"beds"`
	NumberOfReviews      int64        `json:"number_of_reviews"`
	Bathrooms            Bathrooms    `json:"bathrooms"`
	Amenities            []string     `json:"amenities"`
	Price                Bathrooms    `json:"price"`
	SecurityDeposit      Bathrooms    `json:"security_deposit"`
	CleaningFee          Bathrooms    `json:"cleaning_fee"`
	ExtraPeople          Bathrooms    `json:"extra_people"`
	GuestsIncluded       Bathrooms    `json:"guests_included"`
	Images               Images       `json:"images"`
	Host                 Host         `json:"host"`
	Address              Address      `json:"address"`
	Availability         Availability `json:"availability"`
	ReviewScores         ReviewScores `json:"review_scores"`
	Reviews              []Review     `json:"reviews"`
}

type Address struct {
	Street         string   `json:"street"`
	Suburb         string   `json:"suburb"`
	GovernmentArea string   `json:"government_area"`
	Market         string   `json:"market"`
	Country        string   `json:"country"`
	CountryCode    string   `json:"country_code"`
	Location       Location `json:"location"`
}

type Location struct {
	Type            string    `json:"type"`
	Coordinates     []float64 `json:"coordinates"`
	IsLocationExact bool      `json:"is_location_exact"`
}

type Availability struct {
	Availability30  int64 `json:"availability_30"`
	Availability60  int64 `json:"availability_60"`
	Availability90  int64 `json:"availability_90"`
	Availability365 int64 `json:"availability_365"`
}

type Bathrooms struct {
	NumberDecimal string `json:"$numberDecimal"`
}

type Host struct {
	HostID                 string   `json:"host_id"`
	HostURL                string   `json:"host_url"`
	HostName               string   `json:"host_name"`
	HostLocation           string   `json:"host_location"`
	HostAbout              string   `json:"host_about"`
	HostResponseTime       string   `json:"host_response_time"`
	HostThumbnailURL       string   `json:"host_thumbnail_url"`
	HostPictureURL         string   `json:"host_picture_url"`
	HostNeighbourhood      string   `json:"host_neighbourhood"`
	HostResponseRate       int64    `json:"host_response_rate"`
	HostIsSuperhost        bool     `json:"host_is_superhost"`
	HostHasProfilePic      bool     `json:"host_has_profile_pic"`
	HostIdentityVerified   bool     `json:"host_identity_verified"`
	HostListingsCount      int64    `json:"host_listings_count"`
	HostTotalListingsCount int64    `json:"host_total_listings_count"`
	HostVerifications      []string `json:"host_verifications"`
}

type Images struct {
	ThumbnailURL string `json:"thumbnail_url"`
	MediumURL    string `json:"medium_url"`
	PictureURL   string `json:"picture_url"`
	XlPictureURL string `json:"xl_picture_url"`
}

type ReviewScores struct {
	ReviewScoresAccuracy      int64 `json:"review_scores_accuracy"`
	ReviewScoresCleanliness   int64 `json:"review_scores_cleanliness"`
	ReviewScoresCheckin       int64 `json:"review_scores_checkin"`
	ReviewScoresCommunication int64 `json:"review_scores_communication"`
	ReviewScoresLocation      int64 `json:"review_scores_location"`
	ReviewScoresValue         int64 `json:"review_scores_value"`
	ReviewScoresRating        int64 `json:"review_scores_rating"`
}

type Review struct {
	ID           string `json:"_id"`
	Date         string `json:"date"`
	ListingID    string `json:"listing_id"`
	ReviewerID   string `json:"reviewer_id"`
	ReviewerName string `json:"reviewer_name"`
	Comments     string `json:"comments"`
}
