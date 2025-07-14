# Application called Waterlogger


### Objective
A web app suitable for use on mobile or desktop browers to store pool and hot tub water parameters, log additions and adjustments, calculate water balance indices, and a line chart the information over time. Users should be able to Create, Read, Update, Delete records from each of the entities described below through data grids. System should allow exporting the data as a) Excel files (a worksheet for each entity) and Markdown reports (Headings, descriptions and sections for each entity with data tables sorted by date with older records at the top).

### Interface 
App should use a modern appearance. Navigation controls should have a dark navy background. Each field should display the full name of the data requested (not just the database field name), including the units that are stored in the database. Remember to provide output units as well, especially for the calculated indices. Provide hover tips for each input field with lengthy descriptions of the meaning of each field. This information is presented below, but you can also research additional information to present.

### System Requirements
The app should provide basic user authentication that runs from a single executable on Ubuntu linux and not require the use of a separate web server. Let user configure application parameters, including the port number that it listens on, using a configuration file. A single user role is needed. Data should be stored in a Sqlite database in the same directory as the executable. Deployment should simply be to execute the executable, perhaps by setting a linux service. Create the database if it is not present. Track app versions and allow migrations of the data between versions of the database.

### Logical Requirements
1. USERS - Multiple users and each can Create, Read, Edit, Delete all records for all pools. Attributes stored for each user: username, email, password, date and time record was created or edited, user who last created or edited the record. All fields are required with no default values. Username, and email should be unique.
2. POOLS - Multiple pools. Attributes stored for each pool: name, volume in gallons, type (hot tub or pool), system description,  date and time record was created or edited, user who last created or edited the record. Only name is required and no default values. Name should be unique.
3. SAMPLE - Multiple samples per pool. Multiple measurements per sample. Attributes stored for each sample: pool (selected from list), date and time, user taking the sample (selected from list), test kit making the measurements (selected from list), date and time record was created or edited, user who last created or edited the record. All fields are required. 
4. KITS - Test kit and equipment. Attributes: Name, Description, Purchased date, Replenished Date, date and time record was created or edited, user who last created or edited the record. Only name is required.
5. MEASUREMENTS - Measurements of water parameters. Multiple measurements per sample, but only one sample per measurement. 
   a. The following measurements will be stored:
      1. FC - Decimal number. Free Chlorine: This measures the amount of chlorine available to sanitize the water and kill bacteria and algae. Unit: parts per million (ppm) or milligrams per liter (mg/L). Ideal Range: 1.0 - 4.0 ppm, depending on factors like pool usage and environmental conditions. Required. No Default.
      2. TC - Decimal number. Total Chlorine  - Decimal number: This is the sum of free chlorine and combined chlorine (chlorine already used in the sanitation process). Unit: ppm or mg/L. Ideal Range: Ideally the same as free chlorine to minimize combined chlorine, typically within the same range as free chlorine.  Required. No Default.
      3. PH - Decimal number. pH: This measures the acidity or alkalinity of the water. Unit: pH scale (0-14, with 7 being neutral). Ideal Range: 7.4 - 7.6.  Required. No Default.
      4. TA - Decimal number. Total Alkalinity (TA): Measures the water's capacity to resist changes in pH (buffering capacity). Unit: ppm or mg/L mg CaCO3/liter. Ideal Range: 80 - 120 ppm.  Required. No Default.
      5. CH - Decimal number. Calcium Hardness (CH): Measures the concentration of dissolved calcium in the pool water. Unit: ppm or mg/L mg Ca/liter. Ideal Range: 200 - 400 ppm.  Required. No Default.
      6. CYA - Decimal number. Cyanuric Acid (CYA): This stabilizes chlorine, protecting it from UV degradation. Unit: ppm or mg/L. Ideal Range: 30 - 50 ppm. Not required. No Default.
      7. T - Decimal number. t, temperature - º;F *  Required. No Default.
      8. SAL - Decimal number. Salinity: The ideal salinity range for most saltwater pools is between 2,700 and 3,400 ppm (parts per million), with 3,200 ppm being optimal. Maintaining the correct salinity ensures the chlorine generator can produce enough chlorine to sanitize the pool effectively and prevents potential damage to pool equipment.  Not required. No Default.
      9. TDS - Decimal number. Total Dissolved Solids - mg/l; Not required. No Default.
      10. APPEARS - Text.  Notes describing the appearance of the water.  Not required. No Default.
      11. MAINT - Text.  Notes on maintenance completed.  Not required. No Default.
   
6. INDICES - Calculated numbers derived from measurements in a single sample. Each sample has one set of each index. Calculate an index if the required measurements are available to complete the calculation. You can research and suggest any additional calculations to perform that are relevant to pool and hot tub operation.

Use calculations from the following URLs to calculate the Ryznar Stability Index (RSI), Langelier Saturation Index (LSI) for 

* https://waterpy.blogspot.com/search/label/Water%20chemistry
* https://github.com/johnzastrow/WaterPy/blob/master/WaterChemistry.py

Here is python code for calculating the indices

```python

#Saturation pH calcite or calcium carbonate
def WATERCHEM_phscalcium(t,tds,ca,hco3):
    tk = UNIT_temperature_c2k(t) # temperature: oC to K
    mca = ca*0.001/40.08 # Ca2+: mg/l to Mole/l
    mhco3= hco3*0.001/100 # Alkalinity, hco3- : mg/l to Mole/l
    i = 2.5*10**(-5)*tds # ionic strenght, Moles/l
    d = 1 #density
    e = 60954/(tk+116)-68.937 #dielectric constant
    a = 1.825*10**6*d**0.5*(e*tk)**(-1.5) #correction factor
    zca = 2 #ionic charge of C2+
    zhco3 = 1 #ionic charge of HCO3-
    if i <= 0.5:
        lghco3 = -a*zhco3**2*(i**0.5/(1+i**0.5)-0.3*i) # activity coefficient for HCO3-
    if i > 0.5:
        lghco3 = -a*zhco3**2*(i**0.5/(1+i**0.5)) # activity coefficient for HCO3-
    lgca = -a*zca**2*(i**0.5/(1+i**0.5)) # activity coefficient for Ca
    pk2 = 2902.39/tk + 0.02379*tk-6.498
    k2 = 10**(-pk2)
    gamma_d = 10**lgca
    kl2= k2/gamma_d
    pkl2 = math.log10(1/kl2)
    pks = 0.01183*t+8.03
    ks = 1/10**pks
    kls = ks/gamma_d**2
    pkls = math.log10(1/kls)
    pca=math.log10(1/mca)
    phs = pkl2 + pca - pkls - math.log10(2*mhco3)- lghco3
    return(phs);

# Langelier Saturation Index - LSI
def WATERCHEM_lsi(t,ph,tds,ca,hco3):
        phs = WATERCHEM_phscalcium(t,tds,ca,hco3)
        lsi = ph - phs
        return(lsi);

# Ryznar Stability Index - RSI
def WATERCHEM_rsi(t,ph,tds,ca,hco3):
    phs = WATERCHEM_phscalcium(t,tds,ca,hco3)
    ri = 2*phs - ph
    return(ri);
```

And here is an explanation of the calculations from the URLs above. Note that this app is being commissioned by an American, who thinks largely in Imperial units, and the equations below are based on international SI units. Be sure to incorporate the correct unit conversions to use the formulas below.

WATERCHEM_rsi(): Ryznar Stability Index (RSI)
About:
Calculate the Ryznar Stability Index (RSI).
RSI=2×pHs−pH

Function: WATERCHEM_rsi(t,ph,tds,ca,hco3)
Parameters:
t, temperature - ºC;
ph, water pH
tds, Total Dissolved Solids - mg/l;
ca, concentration of calcium hardness, mg Ca/liter;
hco3, concetration of HCO3- (alkalinity), mg CaCO3/liter.
Sample code:
from WaterChemistry  import *

ph = 7.7    # water pH
tds = 300   # mg TDS/liter
ca = 50     # mg Ca/liter
hco3 = 100  # mg CaCo3/liter (Alkalinity)
t = 20      # oC

rsi = WATERCHEM_rsi(t,ph,tds,ca,hco3)

print("RSI:",rsi)
Result:
RSI: 8.128539459578171



WATERCHEM_lsi(): Langelier Saturation Index (LSI)


Calculate the Langelier Saturation Index (LSI).
LSI=pH−pHs


Module: WaterChemistry

Function: WATERCHEM_lsi(t,ph,tds,ca,hco3)
Parameters:
t, temperature - ºC;
ph, water pH
tds, Total Dissolved Solids - mg/l;
ca, concentration of calcium hardness, mg Ca/liter;
hco3, concetration of HCO3- (alkalinity), mg CaCO3/liter.
Sample code:
from WaterChemistry  import *

ph = 7.7    # water pH
tds = 300   # mg TDS/liter
ca = 50     # mg Ca/liter
hco3 = 100  # mg CaCo3/liter (Alkalinity)
t = 20      # oC

lsi = WATERCHEM_lsi(t,ph,tds,ca,hco3)

print("LSI:",lsi)

Result:
LSI: -0.21426972978908498

Publicada por WaterPy
Etiquetas: Lagelier Saturation Index, Water chemistry
WATERCHEM_phscalcium(): Saturation pH for calcium carbonate
About:
Calculate the saturation pH for calcium carbonate.


pHs=pK′2+pCa2+−pK′s−log(2×[Alk])−logγm
\begin{equation}

Module: WaterChemistry

Function: WATERCHEM_phscalcium(t,tds,ca,hco3)
Parameters:
t, temperature - ºC;
tds, Total Dissolved Solids - mg/l;
ca, concentration of calcium hardness, mg Ca/liter;
hco3, concetration of HCO3- (alkalinity), mg CaCO3/liter.
Sample code:
from WaterChemistry  import *

tds = 300   # mg TDS/liter
ca = 50     # mg Ca/liter
hco3 = 100  # mg CaCo3/liter (Alkalinity)
t = 20      # oC

phs = WATERCHEM_phscalcium(t,tds,ca,hco3)

print("phs:",phs)

Result:
phs: 7.914650556552682

Publicada por WaterPy
Etiquetas: Calcium carbonate saturation, pH, Water chemistry
