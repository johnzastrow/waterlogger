package adjustments

import (
	"fmt"
	"math"
	"waterlogger/internal/chemistry"
)

// ChemicalAdjustmentRequest represents input for chemical adjustment calculations
type ChemicalAdjustmentRequest struct {
	// Pool information
	PoolVolume float64 `json:"pool_volume"` // gallons
	PoolType   string  `json:"pool_type"`   // "pool" or "hot_tub"
	
	// Starting values (current readings)
	StartingFC   float64  `json:"starting_fc"`   // Free Chlorine (ppm)
	StartingPH   float64  `json:"starting_ph"`   // pH (0-14 scale)
	StartingTA   float64  `json:"starting_ta"`   // Total Alkalinity (ppm)
	StartingCH   float64  `json:"starting_ch"`   // Calcium Hardness (ppm)
	StartingCYA  *float64 `json:"starting_cya"`  // Cyanuric Acid (ppm)
	StartingTemp float64  `json:"starting_temp"` // Temperature (°F)
	StartingSalt *float64 `json:"starting_salt"` // Salinity (ppm)
	StartingTDS  *float64 `json:"starting_tds"`  // Total Dissolved Solids (mg/l)
	
	// Target values (desired readings)
	TargetFC   float64  `json:"target_fc"`   // Free Chlorine (ppm)
	TargetPH   float64  `json:"target_ph"`   // pH (0-14 scale)
	TargetTA   float64  `json:"target_ta"`   // Total Alkalinity (ppm)
	TargetCH   float64  `json:"target_ch"`   // Calcium Hardness (ppm)
	TargetCYA  *float64 `json:"target_cya"`  // Cyanuric Acid (ppm)
	TargetTemp float64  `json:"target_temp"` // Temperature (°F)
	TargetSalt *float64 `json:"target_salt"` // Salinity (ppm)
	TargetTDS  *float64 `json:"target_tds"`  // Total Dissolved Solids (mg/l)
}

// ChemicalAdjustmentResult represents the calculated chemical adjustments
type ChemicalAdjustmentResult struct {
	// Calculated indices
	StartingLSI float64 `json:"starting_lsi"`
	StartingRSI float64 `json:"starting_rsi"`
	TargetLSI   float64 `json:"target_lsi"`
	TargetRSI   float64 `json:"target_rsi"`
	
	// Chemical recommendations (amounts to add)
	AddMuriaticAcid     float64 `json:"add_muriatic_acid"`     // fl oz
	AddSodiumBisulfate  float64 `json:"add_sodium_bisulfate"`  // oz (weight)
	AddSodaAsh          float64 `json:"add_soda_ash"`          // oz (weight)
	AddBorax            float64 `json:"add_borax"`             // oz (weight)
	AddSodiumBicarbonate float64 `json:"add_sodium_bicarbonate"` // oz (weight)
	AddCalciumChloride  float64 `json:"add_calcium_chloride"`  // oz (weight)
	AddBleach           float64 `json:"add_bleach"`            // fl oz
	AddTrichlor         float64 `json:"add_trichlor"`          // oz (weight)
	AddDichlor          float64 `json:"add_dichlor"`           // oz (weight)
	AddCalHypo          float64 `json:"add_cal_hypo"`          // oz (weight)
	AddSalt             float64 `json:"add_salt"`              // lbs
	
	// Adjustment recommendations and warnings
	Recommendations []string `json:"recommendations"`
	Warnings        []string `json:"warnings"`
	Priority        []string `json:"priority"` // Order of chemical additions
}

// CalculateAdjustments calculates chemical adjustments needed to reach target values
func CalculateAdjustments(req ChemicalAdjustmentRequest) (*ChemicalAdjustmentResult, error) {
	result := &ChemicalAdjustmentResult{
		Recommendations: []string{},
		Warnings:        []string{},
		Priority:        []string{},
	}
	
	// Validate pool volume
	if req.PoolVolume <= 0 {
		return nil, fmt.Errorf("pool volume must be greater than 0")
	}
	
	// Calculate starting and target LSI/RSI
	if err := calculateIndices(req, result); err != nil {
		return nil, fmt.Errorf("failed to calculate indices: %v", err)
	}
	
	// Calculate chemical adjustments in proper order
	calculatePHAdjustments(req, result)
	calculateAlkalinityAdjustments(req, result)
	calculateHardnessAdjustments(req, result)
	calculateChlorineAdjustments(req, result)
	calculateSaltAdjustments(req, result)
	
	// Generate recommendations and priority order
	generateRecommendations(req, result)
	generatePriorityOrder(result)
	
	return result, nil
}

// calculateIndices calculates LSI and RSI for starting and target values
func calculateIndices(req ChemicalAdjustmentRequest, result *ChemicalAdjustmentResult) error {
	// Use default values for missing parameters
	startingTDS := 300.0
	if req.StartingTDS != nil {
		startingTDS = *req.StartingTDS
	}
	
	targetTDS := 300.0
	if req.TargetTDS != nil {
		targetTDS = *req.TargetTDS
	}
	
	// Convert temperature from Fahrenheit to Celsius for calculations
	startingTempC := (req.StartingTemp - 32) * 5 / 9
	targetTempC := (req.TargetTemp - 32) * 5 / 9
	
	// Calculate starting indices
	startingLSI := chemistry.CalculateLSI(startingTempC, req.StartingPH, startingTDS, req.StartingCH, req.StartingTA)
	startingRSI := chemistry.CalculateRSI(startingTempC, req.StartingPH, startingTDS, req.StartingCH, req.StartingTA)
	
	// Calculate target indices
	targetLSI := chemistry.CalculateLSI(targetTempC, req.TargetPH, targetTDS, req.TargetCH, req.TargetTA)
	targetRSI := chemistry.CalculateRSI(targetTempC, req.TargetPH, targetTDS, req.TargetCH, req.TargetTA)
	
	result.StartingLSI = startingLSI
	result.StartingRSI = startingRSI
	result.TargetLSI = targetLSI
	result.TargetRSI = targetRSI
	
	return nil
}

// calculatePHAdjustments calculates pH adjustment chemicals
func calculatePHAdjustments(req ChemicalAdjustmentRequest, result *ChemicalAdjustmentResult) {
	pHChange := req.TargetPH - req.StartingPH
	
	if pHChange > 0.05 { // Need to raise pH
		// Soda Ash (Sodium Carbonate) - primary pH increaser
		// Rule: 6 oz per 10,000 gallons raises pH by ~0.2
		sodaAshOz := (pHChange / 0.2) * (req.PoolVolume / 10000) * 6
		result.AddSodaAsh = sodaAshOz // Keep in oz (weight)
		
		// Borax alternative for smaller increases
		if pHChange < 0.2 {
			// Borax: 20 oz per 10,000 gallons raises pH by ~0.1
			boraxOz := (pHChange / 0.1) * (req.PoolVolume / 10000) * 20
			result.AddBorax = boraxOz // Keep in oz (weight)
		}
		
	} else if pHChange < -0.05 { // Need to lower pH
		// Muriatic Acid (31.45% HCl) - primary pH decreaser
		// Rule: 10 fl oz per 10,000 gallons lowers pH by ~0.1-0.2
		muriaticAcidOz := (math.Abs(pHChange) / 0.15) * (req.PoolVolume / 10000) * 10
		result.AddMuriaticAcid = muriaticAcidOz
		
		// Sodium Bisulfate alternative (dry acid)
		// Rule: 1 lb per 10,000 gallons lowers pH by ~0.2 = 16 oz per 10,000 gallons
		sodiumBisulfateOz := (math.Abs(pHChange) / 0.2) * (req.PoolVolume / 10000) * 16
		result.AddSodiumBisulfate = sodiumBisulfateOz // Now in oz (weight)
	}
}

// calculateAlkalinityAdjustments calculates total alkalinity adjustment chemicals
func calculateAlkalinityAdjustments(req ChemicalAdjustmentRequest, result *ChemicalAdjustmentResult) {
	taChange := req.TargetTA - req.StartingTA
	
	if taChange > 5 { // Need to raise TA
		// Sodium Bicarbonate (Baking Soda)
		// Rule: 1.5 lbs per 10,000 gallons raises TA by 10 ppm = 24 oz per 10,000 gallons
		bicarbOz := (taChange / 10) * (req.PoolVolume / 10000) * 24
		result.AddSodiumBicarbonate = bicarbOz // Now in oz (weight)
		
	} else if taChange < -5 { // Need to lower TA
		// Use muriatic acid (same as pH reduction)
		// Lower pH to 7.0-7.2, then aerate to raise pH back up
		if req.StartingPH > 7.2 {
			// Additional acid needed to lower TA
			additionalAcid := (math.Abs(taChange) / 10) * (req.PoolVolume / 10000) * 8
			result.AddMuriaticAcid += additionalAcid
		}
	}
}

// calculateHardnessAdjustments calculates calcium hardness adjustment chemicals
func calculateHardnessAdjustments(req ChemicalAdjustmentRequest, result *ChemicalAdjustmentResult) {
	chChange := req.TargetCH - req.StartingCH
	
	if chChange > 10 { // Need to raise CH
		// Calcium Chloride
		// Rule: 1.25 lbs per 10,000 gallons raises CH by 10 ppm = 20 oz per 10,000 gallons
		calciumChlorideOz := (chChange / 10) * (req.PoolVolume / 10000) * 20
		result.AddCalciumChloride = calciumChlorideOz // Now in oz (weight)
		
	} else if chChange < -10 { // Need to lower CH
		// Can only be lowered by partial drain and refill
		drainPercent := math.Abs(chChange) / req.StartingCH * 100
		result.Recommendations = append(result.Recommendations,
			fmt.Sprintf("Lower calcium hardness by partial drain/refill (~%.0f%% of pool water)", drainPercent))
	}
}

// calculateChlorineAdjustments calculates chlorine adjustment chemicals
func calculateChlorineAdjustments(req ChemicalAdjustmentRequest, result *ChemicalAdjustmentResult) {
	fcChange := req.TargetFC - req.StartingFC
	
	if fcChange > 0.5 { // Need to raise FC
		// Liquid Bleach (12.5% sodium hypochlorite)
		// Rule: ~13 fl oz per 10,000 gallons raises FC by 1 ppm
		bleachOz := fcChange * (req.PoolVolume / 10000) * 13
		result.AddBleach = bleachOz
		
		// Cal-Hypo alternative (65% available chlorine)
		// Rule: ~2 oz per 10,000 gallons raises FC by 1 ppm
		calHypoOz := fcChange * (req.PoolVolume / 10000) * 2
		result.AddCalHypo = calHypoOz // Keep in oz (weight)
		
		// Trichlor (90% available chlorine, stabilized)
		// Rule: ~1.5 oz per 10,000 gallons raises FC by 1 ppm
		if req.StartingCYA != nil && *req.StartingCYA < 50 {
			trichlorOz := fcChange * (req.PoolVolume / 10000) * 1.5
			result.AddTrichlor = trichlorOz // Keep in oz (weight)
		}
		
		// Dichlor (56% available chlorine, stabilized)
		// Rule: ~2.4 oz per 10,000 gallons raises FC by 1 ppm
		if req.StartingCYA != nil && *req.StartingCYA < 50 {
			dichlorOz := fcChange * (req.PoolVolume / 10000) * 2.4
			result.AddDichlor = dichlorOz // Keep in oz (weight)
		}
		
	} else if fcChange < -0.5 { // Need to lower FC
		// Can be done by stopping chlorination and letting it naturally decay
		// Or using sodium thiosulfate (chlorine neutralizer)
		result.Recommendations = append(result.Recommendations,
			"Lower chlorine by stopping chlorination and allowing natural decay, or use chlorine neutralizer")
	}
}

// calculateSaltAdjustments calculates salt adjustment for saltwater pools
func calculateSaltAdjustments(req ChemicalAdjustmentRequest, result *ChemicalAdjustmentResult) {
	if req.StartingSalt == nil || req.TargetSalt == nil {
		return // Not a saltwater pool
	}
	
	saltChange := *req.TargetSalt - *req.StartingSalt
	
	if saltChange > 100 { // Need to add salt
		// Rule: 8.35 lbs of salt per 1000 gallons raises salinity by 1000 ppm
		saltLbs := (saltChange / 1000) * (req.PoolVolume / 1000) * 8.35
		result.AddSalt = saltLbs
		
	} else if saltChange < -100 { // Need to lower salt
		// Can only be lowered by partial drain and refill
		drainPercent := math.Abs(saltChange) / *req.StartingSalt * 100
		result.Recommendations = append(result.Recommendations,
			fmt.Sprintf("Lower salinity by partial drain/refill (~%.0f%% of pool water)", drainPercent))
	}
}

// generateRecommendations generates specific recommendations based on calculations
func generateRecommendations(req ChemicalAdjustmentRequest, result *ChemicalAdjustmentResult) {
	// LSI/RSI recommendations
	if result.TargetLSI < -0.3 {
		result.Warnings = append(result.Warnings, "Target LSI is low - water may be corrosive")
	} else if result.TargetLSI > 0.3 {
		result.Warnings = append(result.Warnings, "Target LSI is high - scaling may occur")
	}
	
	if result.TargetRSI < 5.5 {
		result.Warnings = append(result.Warnings, "Target RSI is low - scaling tendency")
	} else if result.TargetRSI > 7.0 {
		result.Warnings = append(result.Warnings, "Target RSI is high - corrosive tendency")
	}
	
	// Add safety recommendations
	if result.AddMuriaticAcid > 0 {
		result.Recommendations = append(result.Recommendations,
			"Add muriatic acid slowly in small doses. Wait 4-6 hours between applications.")
	}
	
	if result.AddSodaAsh > 0 || result.AddBorax > 0 {
		result.Recommendations = append(result.Recommendations,
			"Pre-dissolve pH increasers before adding to pool. Distribute evenly.")
	}
	
	if result.AddSodiumBicarbonate > 0 {
		result.Recommendations = append(result.Recommendations,
			"Add baking soda gradually. Can be broadcast directly into pool.")
	}
}

// generatePriorityOrder generates the recommended order of chemical additions
func generatePriorityOrder(result *ChemicalAdjustmentResult) {
	// Based on the order from docs/chem_calcs.md:
	// 1. Metal treatment (if needed)
	// 2. Alkalinity
	// 3. pH  
	// 4. Calcium Hardness
	// 5. Sanitizer (Chlorine)
	// 6. CYA/Stabilizer
	// 7. Salt (if applicable)
	
	if result.AddSodiumBicarbonate > 0 {
		result.Priority = append(result.Priority, "Total Alkalinity (Sodium Bicarbonate)")
	}
	
	if result.AddMuriaticAcid > 0 {
		result.Priority = append(result.Priority, "pH Down (Muriatic Acid)")
	} else if result.AddSodiumBisulfate > 0 {
		result.Priority = append(result.Priority, "pH Down (Sodium Bisulfate)")
	}
	
	if result.AddSodaAsh > 0 {
		result.Priority = append(result.Priority, "pH Up (Soda Ash)")
	} else if result.AddBorax > 0 {
		result.Priority = append(result.Priority, "pH Up (Borax)")
	}
	
	if result.AddCalciumChloride > 0 {
		result.Priority = append(result.Priority, "Calcium Hardness (Calcium Chloride)")
	}
	
	if result.AddBleach > 0 {
		result.Priority = append(result.Priority, "Free Chlorine (Liquid Bleach)")
	} else if result.AddCalHypo > 0 {
		result.Priority = append(result.Priority, "Free Chlorine (Cal-Hypo)")
	} else if result.AddTrichlor > 0 {
		result.Priority = append(result.Priority, "Free Chlorine (Trichlor)")
	} else if result.AddDichlor > 0 {
		result.Priority = append(result.Priority, "Free Chlorine (Dichlor)")
	}
	
	if result.AddSalt > 0 {
		result.Priority = append(result.Priority, "Salt")
	}
}

// GetTargetRanges returns recommended target ranges for pool/hot tub parameters
func GetTargetRanges(poolType string) map[string]map[string]float64 {
	ranges := make(map[string]map[string]float64)
	
	if poolType == "hot_tub" {
		ranges["fc"] = map[string]float64{"min": 1.0, "max": 3.0, "ideal": 2.0}
		ranges["ph"] = map[string]float64{"min": 7.4, "max": 7.6, "ideal": 7.5}
		ranges["ta"] = map[string]float64{"min": 100.0, "max": 150.0, "ideal": 125.0}
		ranges["ch"] = map[string]float64{"min": 175.0, "max": 250.0, "ideal": 200.0}
		ranges["lsi"] = map[string]float64{"min": -0.3, "max": 0.3, "ideal": 0.0}
		ranges["rsi"] = map[string]float64{"min": 5.5, "max": 7.0, "ideal": 6.0}
	} else { // pool
		ranges["fc"] = map[string]float64{"min": 1.0, "max": 3.0, "ideal": 2.0}
		ranges["ph"] = map[string]float64{"min": 7.4, "max": 7.6, "ideal": 7.5}
		ranges["ta"] = map[string]float64{"min": 80.0, "max": 120.0, "ideal": 100.0}
		ranges["ch"] = map[string]float64{"min": 200.0, "max": 400.0, "ideal": 300.0}
		ranges["cya"] = map[string]float64{"min": 30.0, "max": 50.0, "ideal": 40.0}
		ranges["lsi"] = map[string]float64{"min": -0.3, "max": 0.3, "ideal": 0.0}
		ranges["rsi"] = map[string]float64{"min": 5.5, "max": 7.0, "ideal": 6.0}
	}
	
	return ranges
}