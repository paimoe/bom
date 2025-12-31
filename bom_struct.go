package main

import "time"

type Response struct {
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
	TempMaxCel *float64    `json:"temp_max_cel,omitempty"`
	TempMinCel *float64    `json:"temp_min_cel,omitempty"`
	Precip     *Precip     `json:"precip,omitempty"`
	Weather    *Weather    `json:"weather,omitempty"`
	Radiation  *Radiation  `json:"radiation,omitempty"`
}

/* ---------------- Precipitation ---------------- */

type Precip struct {
	Exceeding10PctTotalMM *float64 `json:"exceeding_10percentchance_total_mm,omitempty"`
	Exceeding25PctTotalMM *float64 `json:"exceeding_25percentchance_total_mm,omitempty"`
	Exceeding50PctTotalMM *float64 `json:"exceeding_50percentchance_total_mm,omitempty"`
	Exceeding75PctTotalMM *float64 `json:"exceeding_75percentchance_total_mm,omitempty"`

	AnyProbabilityPercent        *float64 `json:"any_probability_percent,omitempty"`
	AnyRestOfDayProbabilityPct   *float64 `json:"any_restofday_probability_percent,omitempty"`
	TenMMProbabilityPercent     *float64 `json:"10mm_probability_percent,omitempty"`
	TwentyFiveMMProbabilityPct  *float64 `json:"25mm_probability_percent,omitempty"`
}

/* ---------------- Weather ---------------- */

type Weather struct {
	IconCode *int `json:"icon_code,omitempty"`
}

/* ---------------- Radiation ---------------- */

type Radiation struct {
	UVClearSkyMaxCode *float64  `json:"uv_clear_sky_max_code,omitempty"`
	UVPeriodStart    *time.Time `json:"uv_period_start,omitempty"`
	UVPeriodEnd      *time.Time `json:"uv_period_end,omitempty"`
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
