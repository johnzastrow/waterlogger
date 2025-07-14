// Debug utility for better error handling and logging
window.WaterloggerDebug = {
    // Enhanced fetch with debugging
    async fetch(url, options = {}) {
        console.log(`[DEBUG] Fetching ${options.method || 'GET'} ${url}`);
        
        if (options.body) {
            console.log('[DEBUG] Request body:', options.body);
        }
        
        try {
            const response = await fetch(url, options);
            console.log(`[DEBUG] Response status: ${response.status}`);
            
            // Clone response to read it twice
            const responseClone = response.clone();
            
            let data;
            const contentType = response.headers.get('content-type');
            
            if (contentType && contentType.includes('application/json')) {
                try {
                    data = await responseClone.json();
                    console.log('[DEBUG] Response data:', data);
                } catch (jsonError) {
                    console.error('[DEBUG] Failed to parse JSON response:', jsonError);
                    const text = await responseClone.text();
                    console.log('[DEBUG] Response text:', text);
                    data = { error: 'Invalid JSON response', details: text };
                }
            } else {
                data = await responseClone.text();
                console.log('[DEBUG] Response text:', data);
            }
            
            return {
                ok: response.ok,
                status: response.status,
                statusText: response.statusText,
                headers: response.headers,
                data: data,
                originalResponse: response
            };
        } catch (error) {
            console.error('[DEBUG] Network error:', error);
            throw error;
        }
    },
    
    // Format error messages
    formatError(error, data) {
        if (data && data.error) {
            let message = data.error;
            if (data.details) {
                const details = Array.isArray(data.details) ? data.details.join(', ') : data.details;
                message += ': ' + details;
            }
            return message;
        }
        return error.message || 'Unknown error occurred';
    },
    
    // Common error handler
    handleError(error, context = '') {
        console.error(`[DEBUG] Error in ${context}:`, error);
        
        if (error.name === 'TypeError' && error.message.includes('fetch')) {
            return 'Network error: Unable to connect to server';
        }
        
        return this.formatError(error);
    }
};

// Common Alpine.js helpers
window.WaterloggerHelpers = {
    // Standard form submission handler
    async submitForm(formData, url, method = 'POST', context = '') {
        console.log(`[DEBUG] Submitting form to ${url}:`, formData);
        
        try {
            const result = await WaterloggerDebug.fetch(url, {
                method: method,
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData)
            });
            
            if (result.ok) {
                console.log(`[DEBUG] ${context} successful`);
                return { success: true, data: result.data };
            } else {
                const errorMessage = WaterloggerDebug.formatError(null, result.data);
                console.error(`[DEBUG] ${context} failed:`, errorMessage);
                return { success: false, error: errorMessage };
            }
        } catch (error) {
            const errorMessage = WaterloggerDebug.handleError(error, context);
            return { success: false, error: errorMessage };
        }
    },
    
    // Standard data loading handler
    async loadData(url, context = '') {
        console.log(`[DEBUG] Loading data from ${url}`);
        
        try {
            const result = await WaterloggerDebug.fetch(url);
            
            if (result.ok) {
                console.log(`[DEBUG] ${context} loaded successfully:`, result.data);
                return { success: true, data: result.data };
            } else {
                const errorMessage = WaterloggerDebug.formatError(null, result.data);
                console.error(`[DEBUG] Failed to load ${context}:`, errorMessage);
                return { success: false, error: errorMessage };
            }
        } catch (error) {
            const errorMessage = WaterloggerDebug.handleError(error, `loading ${context}`);
            return { success: false, error: errorMessage };
        }
    },
    
    // Unit conversion helper
    async convertUnits(value, parameter, fromSystem) {
        console.log(`[DEBUG] Converting ${value} ${parameter} from ${fromSystem}`);
        
        try {
            const result = await WaterloggerDebug.fetch('/api/convert', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    value: value,
                    parameter: parameter,
                    from_system: fromSystem
                })
            });
            
            if (result.ok) {
                console.log(`[DEBUG] Unit conversion successful:`, result.data);
                return { success: true, data: result.data };
            } else {
                const errorMessage = WaterloggerDebug.formatError(null, result.data);
                console.error(`[DEBUG] Unit conversion failed:`, errorMessage);
                return { success: false, error: errorMessage };
            }
        } catch (error) {
            const errorMessage = WaterloggerDebug.handleError(error, 'unit conversion');
            return { success: false, error: errorMessage };
        }
    }
};

// Unit conversion utilities
window.WaterloggerUnits = {
    // Format a measurement value with dual-unit display
    formatMeasurement(value, parameter, userSystem = 'imperial') {
        if (!value || value === 0) return '';
        
        const conversions = {
            temperature: {
                imperial: { unit: '°F', convert: (f) => ((f - 32) * 5/9), convertedUnit: '°C' },
                metric: { unit: '°C', convert: (c) => (c * 9/5 + 32), convertedUnit: '°F' }
            },
            volume: {
                imperial: { unit: 'gal', convert: (gal) => (gal * 3.78541), convertedUnit: 'L' },
                metric: { unit: 'L', convert: (l) => (l / 3.78541), convertedUnit: 'gal' }
            }
        };
        
        if (conversions[parameter]) {
            const config = conversions[parameter][userSystem];
            const converted = config.convert(value);
            return `${value} ${config.unit} (${converted.toFixed(1)} ${config.convertedUnit})`;
        }
        
        // For parameters without conversion (pH, ppm values)
        const units = {
            ph: '',
            fc: ' ppm',
            tc: ' ppm',
            ta: ' ppm',
            ch: ' ppm',
            cya: ' ppm',
            salinity: ' ppm',
            tds: ' ppm',
            lsi: '',
            rsi: '',
            csi: ''
        };
        
        return `${value}${units[parameter] || ''}`;
    },
    
    // Get the appropriate unit label for a parameter
    getUnitLabel(parameter, userSystem = 'imperial') {
        const units = {
            temperature: userSystem === 'imperial' ? '°F' : '°C',
            volume: userSystem === 'imperial' ? 'gal' : 'L',
            ph: '',
            fc: 'ppm',
            tc: 'ppm',
            ta: 'ppm',
            ch: 'ppm',
            cya: 'ppm',
            salinity: 'ppm',
            tds: 'ppm',
            lsi: '',
            rsi: '',
            csi: ''
        };
        
        return units[parameter] || '';
    }
};