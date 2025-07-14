package chemistry

import (
	"fmt"
	"math"
	"strings"
	"waterlogger/internal/models"
)

// Mid-range defaults for missing parameters
const (
	DefaultTDS           = 300.0  // mg/L
	DefaultCalciumHardness = 250.0  // ppm
	DefaultTotalAlkalinity = 100.0  // ppm
)

// Unit conversion functions
func FahrenheitToCelsius(fahrenheit float64) float64 {
	return (fahrenheit - 32) * 5.0 / 9.0
}

func CelsiusToFahrenheit(celsius float64) float64 {
	return celsius*9.0/5.0 + 32
}

// CalculatePhSCalcium calculates the saturation pH for calcium carbonate
func CalculatePhSCalcium(tempC, tds, ca, hco3 float64) float64 {
	tk := tempC + 273.15 // temperature: °C to K
	mca := ca * 0.001 / 40.08 // Ca2+: mg/l to Mole/l
	mhco3 := hco3 * 0.001 / 100 // Alkalinity, hco3-: mg/l to Mole/l
	i := 2.5e-5 * tds // ionic strength, Moles/l
	d := 1.0 // density
	e := 60954/(tk+116) - 68.937 // dielectric constant
	a := 1.825e6 * math.Pow(d, 0.5) * math.Pow(e*tk, -1.5) // correction factor
	zca := 2.0 // ionic charge of Ca2+
	zhco3 := 1.0 // ionic charge of HCO3-
	
	var lghco3 float64
	if i <= 0.5 {
		lghco3 = -a * math.Pow(zhco3, 2) * (math.Pow(i, 0.5)/(1+math.Pow(i, 0.5)) - 0.3*i)
	} else {
		lghco3 = -a * math.Pow(zhco3, 2) * (math.Pow(i, 0.5) / (1 + math.Pow(i, 0.5)))
	}
	
	lgca := -a * math.Pow(zca, 2) * (math.Pow(i, 0.5) / (1 + math.Pow(i, 0.5)))
	pk2 := 2902.39/tk + 0.02379*tk - 6.498
	k2 := math.Pow(10, -pk2)
	gammaD := math.Pow(10, lgca)
	kl2 := k2 / gammaD
	pkl2 := math.Log10(1 / kl2)
	pks := 0.01183*tempC + 8.03
	ks := 1 / math.Pow(10, pks)
	kls := ks / math.Pow(gammaD, 2)
	pkls := math.Log10(1 / kls)
	pca := math.Log10(1 / mca)
	
	phs := pkl2 + pca - pkls - math.Log10(2*mhco3) - lghco3
	return phs
}

// CalculateLSI calculates the Langelier Saturation Index
func CalculateLSI(tempC, ph, tds, ca, hco3 float64) float64 {
	phs := CalculatePhSCalcium(tempC, tds, ca, hco3)
	return ph - phs
}

// CalculateRSI calculates the Ryznar Stability Index
func CalculateRSI(tempC, ph, tds, ca, hco3 float64) float64 {
	phs := CalculatePhSCalcium(tempC, tds, ca, hco3)
	return 2*phs - ph
}

// CalculateIndices calculates LSI and RSI for a measurement with default handling
func CalculateIndices(m *models.Measurements) (*models.Indices, error) {
	if m == nil {
		return nil, fmt.Errorf("measurements cannot be nil")
	}

	// Convert temperature from Fahrenheit to Celsius
	tempC := FahrenheitToCelsius(m.Temperature)
	
	// Use defaults for missing parameters and track them
	var missingParams []string
	
	tds := DefaultTDS
	if m.TDS != nil {
		tds = *m.TDS
	} else {
		missingParams = append(missingParams, "TDS")
	}
	
	ca := DefaultCalciumHardness
	if m.CH != 0 {
		ca = m.CH
	} else {
		missingParams = append(missingParams, "Calcium Hardness")
	}
	
	hco3 := DefaultTotalAlkalinity
	if m.TA != 0 {
		hco3 = m.TA
	} else {
		missingParams = append(missingParams, "Total Alkalinity")
	}
	
	// Calculate indices
	lsi := CalculateLSI(tempC, m.PH, tds, ca, hco3)
	rsi := CalculateRSI(tempC, m.PH, tds, ca, hco3)
	
	// Create comment if parameters were estimated
	var comment *string
	if len(missingParams) > 0 {
		commentText := fmt.Sprintf("Estimated. Calculated with mid-range defaults for the following parameters that were missing: %s", strings.Join(missingParams, ", "))
		comment = &commentText
	}
	
	return &models.Indices{
		BaseModel: models.BaseModel{},
		SampleID:  m.SampleID,
		LSI:       &lsi,
		RSI:       &rsi,
		Comment:   comment,
	}, nil
}

// GetIdealRanges returns ideal ranges for water parameters
func GetIdealRanges() map[string]string {
	return map[string]string{
		"fc":          "1.0 - 4.0 ppm",
		"tc":          "Same as FC (minimize combined chlorine)",
		"ph":          "7.4 - 7.6",
		"ta":          "80 - 120 ppm",
		"ch":          "200 - 400 ppm",
		"cya":         "30 - 50 ppm",
		"salinity":    "2,700 - 3,400 ppm (optimal: 3,200 ppm)",
		"lsi":         "-0.3 to +0.3 (balanced water)",
		"rsi":         "6.0 - 7.0 (stable water)",
	}
}

// GetParameterDescriptions returns detailed descriptions for tooltips
func GetParameterDescriptions() map[string]string {
	return map[string]string{
		"fc": "Free Chlorine measures the amount of chlorine available to sanitize the water and kill bacteria and algae. This is the active form of chlorine that provides ongoing protection.",
		"tc": "Total Chlorine is the sum of free chlorine and combined chlorine (chlorine already used in the sanitation process). Ideally, this should be close to free chlorine levels.",
		"ph": "pH measures the acidity or alkalinity of the water on a scale from 0-14, with 7 being neutral. Proper pH is crucial for chlorine effectiveness and swimmer comfort.",
		"ta": "Total Alkalinity measures the water's capacity to resist changes in pH (buffering capacity). It helps stabilize pH levels and prevents rapid pH swings.",
		"ch": "Calcium Hardness measures the concentration of dissolved calcium in the pool water. Proper levels prevent water from becoming corrosive or causing scale formation.",
		"cya": "Cyanuric Acid stabilizes chlorine, protecting it from UV degradation. It acts as a sunscreen for chlorine but can reduce its effectiveness at high levels.",
		"temperature": "Water temperature affects chemical reaction rates, chlorine effectiveness, and swimmer comfort. Higher temperatures require more sanitizer.",
		"salinity": "Salinity measures dissolved salt content in saltwater pools. Proper levels ensure the chlorine generator can produce adequate chlorine for sanitation.",
		"tds": "Total Dissolved Solids measures all dissolved substances in the water. High TDS can interfere with chemical effectiveness and water clarity.",
		"appearance": "Visual observations about water clarity, color, or any visible issues that may indicate water quality problems.",
		"maintenance": "Notes about maintenance activities performed, equipment issues, or other relevant information about pool care.",
		"lsi": "Langelier Saturation Index indicates whether water is balanced, scale-forming, or corrosive. Values near zero indicate balanced water.",
		"rsi": "Ryznar Stability Index predicts the tendency of water to precipitate or dissolve calcium carbonate. Lower values indicate scale-forming tendency.",
	}
}