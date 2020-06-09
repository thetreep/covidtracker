package cds

import "github.com/thetreep/covidtracker"

type AuthentificationReq struct {
	Username string `json:"Username,omitempty"`
	Password string `json:"Password,omitempty"`
}

type AuthentificationResp struct {
	APIHTTPStatusCode int    `json:"ApiHttpStatusCode,omitempty"`
	Message           string `json:"Message,omitempty"`
	RequestedURI      string `json:"RequestedUri,omitempty"`
	Token             string `json:"Token,omitempty"`
	AgentDutyCode     string `json:"agentDutyCode,omitempty"`
	Login             string `json:"login,omitempty"`
}

type hotel struct {
	HtlID         int     `json:"HtlId"`
	HtlCd         string  `json:"HtlCd"`
	HtlName       string  `json:"HtlName"`
	HtlAddress1   string  `json:"HtlAddress1"`
	HtlCity       string  `json:"HtlCity"`
	HtlZipCode    string  `json:"HtlZipCode"`
	CntCd         string  `json:"CntCd"`
	HtlDesc1      string  `json:"HtlDesc1"`
	HtlDesc2      string  `json:"HtlDesc2"`
	HtlStars      int     `json:"HtlStars"`
	HtlLongitude  float64 `json:"HtlLongitude"`
	HtlLatitude   float64 `json:"HtlLatitude"`
	HtlPhone      string  `json:"HtlPhone"`
	HtlFax        string  `json:"HtlFax"`
	HtlEmail      string  `json:"HtlEmail"`
	HtlNameSearch string  `json:"HtlNameSearch"`
	ImageURL      string  `json:"ImageUrl"`
	HotelInfoList []struct {
		OtaCd   int64 `json:"OtaCd"`
		OtaCode struct {
			OtaName  string `json:"OtaName"`
			OtaName2 string `json:"OtaName2"`
		} `json:"OtaCode"`
		OtaType string `json:"OtaType"`
	} `json:"HotelInfoList"`
	RoomAmenityList []struct {
		OtaCd   int64 `json:"OtaCd"`
		OtaCode struct {
			OtaName  string `json:"OtaName"`
			OtaName2 string `json:"OtaName2"`
		} `json:"OtaCode"`
		OtaType string `json:"OtaType"`
	} `json:"RoomAmenityList"`
	ImageInfoList []struct {
		ImgFile string `json:"ImgFile"`
	} `json:"ImageInfoList"`
	AditionalInfoDescriptionListFr []string    `json:"AditionalInfoDescriptionListFr"`
	AditionalInfoDescriptionListEn []string    `json:"AditionalInfoDescriptionListEn"`
	SanitaryNote                   float64     `json:"SanitaryNote"`
	SanitaryNorm                   string      `json:"SanitaryNorm"`
	ScoreBookingCom                string      `json:"ScoreBookingCom"`
	NbrReviewBookingCom            string      `json:"NbrReviewBookingCom"`
	ScoreTripAdvisor               string      `json:"ScoreTripAdvisor"`
	NbrReviewTripAdvisor           string      `json:"NbrReviewTripAdvisor"`
	ImgTripAdvisor                 string      `json:"ImgTripAdvisor"`
	IsPhoneBookingActivated        interface{} `json:"IsPhoneBookingActivated"`
	IsCvvRequired                  bool        `json:"IsCvvRequired"`
	CdsHotelLink                   string      `json:"CdsHotelLink"`
}

type hotelResultResp struct {
	APIHTTPStatusCode int     `json:"ApiHttpStatusCode"`
	Hotels            []hotel `json:"Hotels"`
	Message           string  `json:"Message"`
	RequestedURI      string  `json:"RequestedUri"`
	Result            int     `json:"Result"`
}

func (h hotel) ToHotel() *covidtracker.Hotel {
	hotel := &covidtracker.Hotel{
		Name:          h.HtlName,
		Address:       h.HtlAddress1,
		City:          h.HtlCity,
		ZipCode:       h.HtlZipCode,
		Country:       h.CntCd,
		ImageURL:      h.ImageURL,
		SanitaryInfos: h.AditionalInfoDescriptionListFr,
		SanitaryNote:  h.SanitaryNote,
		SanitaryNorm:  h.SanitaryNorm,
	}
	return hotel
}
