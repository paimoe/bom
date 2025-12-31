package main

import "time"

type ResponseFcast struct {
	Meta Meta `json:"meta"`
	Fcst Fcst `json:"fcst"`
}

/* ---------------- Meta ---------------- */

type Meta struct {
	IssueTimeUTC     time.Time `json:"issue_time_utc"`
	IssueTimeNextUTC time.Time `json:"issue_time_next_utc"`
	LocalTimezone    string    `json:"local_timezone"`
}

/* ---------------- Forecast Root ---------------- */

type Fcst struct {
	Daily []DailyForecast `json:"daily"`
}

/* ---------------- Daily Forecast ---------------- */

type DailyForecast struct {
	DateUTC time.Time `json:"date_utc"`

	Atm   *Atmosphere `json:"atm,omitempty"`
	Terr  *Terrain    `json:"terr,omitempty"`
	Ocn   *Ocean      `json:"ocn,omitempty"`
	Astro *Astro      `json:"astro,omitempty"`
}

/* ---------------- Atmosphere ---------------- */

type Atmosphere struct {
	SurfAir *SurfaceAir `json:"surf_air,omitempty"`
}

type SurfaceAir struct {
	TempMaxCel *float64   `json:"temp_max_cel,omitempty"`
	TempMinCel *float64   `json:"temp_min_cel,omitempty"`
	Precip     *Precip    `json:"precip,omitempty"`
	Weather    *Weather   `json:"weather,omitempty"`
	Radiation  *Radiation `json:"radiation,omitempty"`
}

/* ---------------- Precipitation ---------------- */

type Precip struct {
	Exceeding10PctTotalMM *float64 `json:"exceeding_10percentchance_total_mm,omitempty"`
	Exceeding25PctTotalMM *float64 `json:"exceeding_25percentchance_total_mm,omitempty"`
	Exceeding50PctTotalMM *float64 `json:"exceeding_50percentchance_total_mm,omitempty"`
	Exceeding75PctTotalMM *float64 `json:"exceeding_75percentchance_total_mm,omitempty"`

	AnyProbabilityPercent      *float64 `json:"any_probability_percent,omitempty"`
	AnyRestOfDayProbabilityPct *float64 `json:"any_restofday_probability_percent,omitempty"`
	TenMMProbabilityPercent    *float64 `json:"10mm_probability_percent,omitempty"`
	TwentyFiveMMProbabilityPct *float64 `json:"25mm_probability_percent,omitempty"`
}

/* ---------------- Weather ---------------- */

type Weather struct {
	IconCode *int `json:"icon_code,omitempty"`
}

/* ---------------- Radiation ---------------- */

type Radiation struct {
	UVClearSkyMaxCode *float64   `json:"uv_clear_sky_max_code,omitempty"`
	UVPeriodStart     *time.Time `json:"uv_period_start,omitempty"`
	UVPeriodEnd       *time.Time `json:"uv_period_end,omitempty"`
}

/* ---------------- Terrain / Ocean / Astro ---------------- */

type Terrain struct {
	SurfLand *SurfaceLand `json:"surf_land,omitempty"`
}

type SurfaceLand struct {
	Snow map[string]any `json:"snow,omitempty"`
}

type Ocean struct {
	SurfWater *SurfaceWater `json:"surf_water,omitempty"`
}

type SurfaceWater struct {
	Sea map[string]any `json:"sea,omitempty"`
}

type Astro struct {
	SunriseUTC *time.Time `json:"sunrise_utc"`
	SunsetUTC  *time.Time `json:"sunset_utc"`
}

type ResponseNow struct {
	Stn Station `json:"stn"`
	Obs Obs     `json:"obs"`
}

/* ---------------- Station ---------------- */

type Station struct {
	Identity StationIdentity `json:"identity"`
	Location StationLocation `json:"location"`
}

type StationIdentity struct {
	BOMStnNum   int     `json:"bom_stn_num"`
	RiverStnID  *string `json:"river_stn_id"` // null in sample; type unknown, keep as *string
	BOMStnName  string  `json:"bom_stn_name"`
	WMOStnID    int     `json:"wmo_stn_id"`
	WIGOSStnID  *string `json:"wigos_stn_id"` // null in sample; type unknown, keep as *string
	HtAboveMSL  float64 `json:"ht_above_msl"`
	HtBarometer float64 `json:"ht_barometer"`
}

type StationLocation struct {
	LatDecDeg  float64 `json:"lat_dec_deg"`
	LongDecDeg float64 `json:"long_dec_deg"`
	Timezone   string  `json:"timezone"`
}

/* ---------------- Observations ---------------- */

type Obs struct {
	DatetimeUTC time.Time `json:"datetime_utc"`

	Temp       Temp       `json:"temp"`
	Pres       Pressure   `json:"pres"`
	Wind       Wind       `json:"wind"`
	Precip     PrecipObs  `json:"precip"`
	Visibility Visibility `json:"visibility"`
	Cloud      Cloud      `json:"cloud"`
}

/* ---------------- Temp ---------------- */

type Temp struct {
	DryBulb1MinCel       float64   `json:"dry_bulb_1min_cel"`
	Apparent1MinCel      float64   `json:"apparent_1min_cel"`
	DewPnt1MinCel        float64   `json:"dew_pnt_1min_cel"`
	WetBulb1MinAvgCel    float64   `json:"wet_bulb_1min_avg_cel"`
	WetBulbGlobeSunCel   float64   `json:"wet_bulb_globe_sun_cel"`
	WetBulbGlobeShadeCel float64   `json:"wet_bulb_globe_shade_cel"`
	WetBulbDepressionCel float64   `json:"wet_bulb_depression_cel"`
	DryBulbMaxCel        float64   `json:"dry_bulb_max_cel"`
	DryBulbMaxTimeUTC    time.Time `json:"dry_bulb_max_time_utc"`
	DryBulbMinCel        float64   `json:"dry_bulb_min_cel"`
	DryBulbMinTimeUTC    time.Time `json:"dry_bulb_min_time_utc"`
	RelHumPercent        float64   `json:"rel_hum_percent"`
}

/* ---------------- Pressure ---------------- */

type Pressure struct {
	StnLvlHPA *float64 `json:"stn_lvl_hpa"` // null in sample
	MSLHPA    float64  `json:"msl_hpa"`
	QNHHPA    float64  `json:"qnh_hpa"`
}

/* ---------------- Wind ---------------- */

type Wind struct {
	Speed10MMPS        float64   `json:"speed_10m_mps"`
	Dirn10MOrd         string    `json:"dirn_10m_ord"`
	GustSpeed10MMPS    float64   `json:"gust_speed_10m_mps"`
	GustDirn10MDegT    float64   `json:"gust_dirn_10m_deg_t"`
	GustSpeed10MMaxMPS float64   `json:"gust_speed_10m_max_mps"`
	Gust10MMaxUTC      time.Time `json:"gust_10m_max_utc"`
	Run2MTotalM        *float64  `json:"run_2m_total_m"` // null in sample
}

/* ---------------- Precip ---------------- */

type PrecipObs struct {
	Since0900LctTotalMM float64  `json:"since_0900lct_total_mm"`
	Since0000LctTotalMM float64  `json:"since_0000lct_total_mm"`
	H24_0900LctTotalMM  float64  `json:"24h_0900lct_total_mm"`
	Min10TotalMM        float64  `json:"10min_total_mm"`
	H1TotalMM           float64  `json:"1h_total_mm"`
	H24TotalMM          *float64 `json:"24h_total_mm"` // null in sample
}

/* ---------------- Visibility ---------------- */

type Visibility struct {
	HorizM *float64 `json:"horiz_m"` // null in sample
}

/* ---------------- Cloud ---------------- */

type Cloud struct {
	BaseHtS1M *float64 `json:"base_ht_s1_m"`
	BaseHtS2M *float64 `json:"base_ht_s2_m"`
	BaseHtS3M *float64 `json:"base_ht_s3_m"`
	BaseHtS4M *float64 `json:"base_ht_s4_m"`
	BaseHtS5M *float64 `json:"base_ht_s5_m"`

	TotalCoverAmtText     *string  `json:"total_cover_amt_text"`
	LowLayerCoverAmtOkta  *float64 `json:"low_layer_cover_amt_okta"`
	LowLayerHeightM       *float64 `json:"low_layer_height_m"`
	MedLayerCoverAmtOkta  *float64 `json:"med_layer_cover_amt_okta"`
	MedLayerHeightM       *float64 `json:"med_layer_height_m"`
	HighLayerCoverAmtOkta *float64 `json:"high_layer_cover_amt_okta"`
	HighLayerHeightM      *float64 `json:"high_layer_height_m"`
}
