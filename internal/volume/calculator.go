package volume

import (
	"fmt"
	"math"
)

// VolumeCalculationRequest represents input for volume calculation
type VolumeCalculationRequest struct {
	Shape        string  `json:"shape"`         // rectangular, round, oval, kidney, irregular
	Length       float64 `json:"length"`        // inches
	Width        float64 `json:"width"`         // inches
	ShallowDepth float64 `json:"shallow_depth"` // inches
	DeepDepth    float64 `json:"deep_depth"`    // inches
	KidneyWidthA float64 `json:"kidney_width_a,omitempty"` // inches (for kidney shape)
	KidneyWidthB float64 `json:"kidney_width_b,omitempty"` // inches (for kidney shape)
}

// VolumeCalculationResult represents the result of volume calculation
type VolumeCalculationResult struct {
	VolumeGallons float64 `json:"volume_gallons"`
	Formula       string  `json:"formula"`
	Explanation   string  `json:"explanation"`
}

// CalculateVolume calculates pool volume based on shape and dimensions
func CalculateVolume(req VolumeCalculationRequest) (*VolumeCalculationResult, error) {
	if req.Shape == "" {
		return nil, fmt.Errorf("shape is required")
	}

	// Convert dimensions from inches to feet for calculations
	lengthFt := req.Length / 12.0
	widthFt := req.Width / 12.0
	shallowDepthFt := req.ShallowDepth / 12.0
	deepDepthFt := req.DeepDepth / 12.0
	
	// Calculate average depth
	avgDepthFt := (shallowDepthFt + deepDepthFt) / 2.0
	
	var volumeCubicFeet float64
	var formula string
	var explanation string

	switch req.Shape {
	case "rectangular":
		// Length × Width × Average Depth × 7.5 = Volume (gallons)
		volumeCubicFeet = lengthFt * widthFt * avgDepthFt
		formula = "Length × Width × Average Depth × 7.5"
		explanation = fmt.Sprintf("%.1f ft × %.1f ft × %.1f ft × 7.5 = %.0f gallons", 
			lengthFt, widthFt, avgDepthFt, volumeCubicFeet * 7.5)

	case "round":
		// 3.14 (π) × Radius × Radius × Average Depth × 7.5 = Volume (gallons)
		// For round pools, we use width as diameter, so radius = width/2
		radiusFt := widthFt / 2.0
		volumeCubicFeet = math.Pi * radiusFt * radiusFt * avgDepthFt
		formula = "π × Radius² × Average Depth × 7.5"
		explanation = fmt.Sprintf("π × %.1f² ft × %.1f ft × 7.5 = %.0f gallons", 
			radiusFt, avgDepthFt, volumeCubicFeet * 7.5)

	case "oval":
		// 3.14 (π) × Length × Width × 0.25 × Average Depth × 7.5 = Volume (gallons)
		volumeCubicFeet = math.Pi * lengthFt * widthFt * 0.25 * avgDepthFt
		formula = "π × Length × Width × 0.25 × Average Depth × 7.5"
		explanation = fmt.Sprintf("π × %.1f ft × %.1f ft × 0.25 × %.1f ft × 7.5 = %.0f gallons", 
			lengthFt, widthFt, avgDepthFt, volumeCubicFeet * 7.5)

	case "kidney":
		// (A + B) × Length × 0.45 × Average Depth × 7.5 = Volume (gallons)
		if req.KidneyWidthA <= 0 || req.KidneyWidthB <= 0 {
			return nil, fmt.Errorf("kidney shape requires both width A and width B")
		}
		kidneyWidthAFt := req.KidneyWidthA / 12.0
		kidneyWidthBFt := req.KidneyWidthB / 12.0
		volumeCubicFeet = (kidneyWidthAFt + kidneyWidthBFt) * lengthFt * 0.45 * avgDepthFt
		formula = "(Width A + Width B) × Length × 0.45 × Average Depth × 7.5"
		explanation = fmt.Sprintf("(%.1f ft + %.1f ft) × %.1f ft × 0.45 × %.1f ft × 7.5 = %.0f gallons", 
			kidneyWidthAFt, kidneyWidthBFt, lengthFt, avgDepthFt, volumeCubicFeet * 7.5)

	case "irregular":
		// Longest Length × Widest Width × Average Depth × 5.9 = Volume (gallons)
		volumeCubicFeet = lengthFt * widthFt * avgDepthFt * (5.9 / 7.5) // Adjust for 5.9 multiplier
		formula = "Longest Length × Widest Width × Average Depth × 5.9"
		explanation = fmt.Sprintf("%.1f ft × %.1f ft × %.1f ft × 5.9 = %.0f gallons", 
			lengthFt, widthFt, avgDepthFt, volumeCubicFeet * 7.5 * (5.9 / 7.5))
		// For irregular, we use 5.9 instead of 7.5
		return &VolumeCalculationResult{
			VolumeGallons: volumeCubicFeet * 7.5 * (5.9 / 7.5),
			Formula:       formula,
			Explanation:   explanation,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported shape: %s", req.Shape)
	}

	// Convert cubic feet to gallons (1 cubic foot = 7.5 gallons)
	volumeGallons := volumeCubicFeet * 7.5

	return &VolumeCalculationResult{
		VolumeGallons: volumeGallons,
		Formula:       formula,
		Explanation:   explanation,
	}, nil
}

// ValidateVolumeRequest validates the volume calculation request
func ValidateVolumeRequest(req VolumeCalculationRequest) error {
	if req.Shape == "" {
		return fmt.Errorf("shape is required")
	}

	validShapes := map[string]bool{
		"rectangular": true,
		"round":       true,
		"oval":        true,
		"kidney":      true,
		"irregular":   true,
	}

	if !validShapes[req.Shape] {
		return fmt.Errorf("invalid shape: %s", req.Shape)
	}

	// Validate length (not required for round pools)
	if req.Shape != "round" && req.Length <= 0 {
		return fmt.Errorf("length must be greater than 0")
	}

	if req.Width <= 0 {
		return fmt.Errorf("width must be greater than 0")
	}

	if req.ShallowDepth <= 0 {
		return fmt.Errorf("shallow depth must be greater than 0")
	}

	if req.DeepDepth <= 0 {
		return fmt.Errorf("deep depth must be greater than 0")
	}

	if req.DeepDepth < req.ShallowDepth {
		return fmt.Errorf("deep depth must be greater than or equal to shallow depth")
	}

	// Special validation for kidney shape
	if req.Shape == "kidney" {
		if req.KidneyWidthA <= 0 {
			return fmt.Errorf("kidney width A must be greater than 0")
		}
		if req.KidneyWidthB <= 0 {
			return fmt.Errorf("kidney width B must be greater than 0")
		}
	}

	return nil
}

// GetShapeRequiredFields returns the required fields for each shape
func GetShapeRequiredFields(shape string) []string {
	baseFields := []string{"length", "width", "shallow_depth", "deep_depth"}
	
	switch shape {
	case "kidney":
		return append(baseFields, "kidney_width_a", "kidney_width_b")
	case "round":
		// For round pools, width represents diameter
		return []string{"width", "shallow_depth", "deep_depth"}
	default:
		return baseFields
	}
}

// GetShapeDescription returns a description of what each dimension means for a shape
func GetShapeDescription(shape string) map[string]string {
	descriptions := make(map[string]string)
	
	switch shape {
	case "rectangular":
		descriptions["length"] = "Length of the pool (inches)"
		descriptions["width"] = "Width of the pool (inches)"
		descriptions["shallow_depth"] = "Depth at shallow end (inches)"
		descriptions["deep_depth"] = "Depth at deep end (inches)"
		
	case "round":
		descriptions["width"] = "Diameter of the pool (inches)"
		descriptions["shallow_depth"] = "Depth at shallow end (inches)"
		descriptions["deep_depth"] = "Depth at deep end (inches)"
		
	case "oval":
		descriptions["length"] = "Length of the oval (inches)"
		descriptions["width"] = "Width of the oval (inches)"
		descriptions["shallow_depth"] = "Depth at shallow end (inches)"
		descriptions["deep_depth"] = "Depth at deep end (inches)"
		
	case "kidney":
		descriptions["length"] = "Length of the kidney shape (inches)"
		descriptions["kidney_width_a"] = "First width measurement (widest point A, inches)"
		descriptions["kidney_width_b"] = "Second width measurement (widest point B, inches)"
		descriptions["shallow_depth"] = "Depth at shallow end (inches)"
		descriptions["deep_depth"] = "Depth at deep end (inches)"
		
	case "irregular":
		descriptions["length"] = "Longest length dimension (inches)"
		descriptions["width"] = "Widest width dimension (inches)"
		descriptions["shallow_depth"] = "Depth at shallow end (inches)"
		descriptions["deep_depth"] = "Depth at deep end (inches)"
	}
	
	return descriptions
}