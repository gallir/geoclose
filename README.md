# geoclose

## Usage
`./geoclose -d giata_hotels.csv  -s toSearch.csv  -o out.csv`

```
$ ./geoclose -h
Usage of ./geoclose:
  -d string
    	Data file, for example giata.csv
  -o string
    	Output CSV filename, if not specified, it prints in stdout
  -s string
    	Data to look for example new_properties.csv
```

## CSV format
Mandatory columns:

`id (int), latitude (float), longitude (float)`

Addtional columns may be added, their values will be summarized in the output.

### CSV Example

#### To search (or needles)
```
    id,latitude,longitude
    10043289,52.847035,1.477516
    10044206,24.448604,124.167447
    10096121,44.659429,-1.1737
    10104857,49.391318,-119.906696
    10104962,49.391012,-119.904283
    10217263,44.099415,9.738465
    10227010,49.720348,-118.922762
    10260297,-40.142316,-71.296056
```

#### Data (or haystack)
```
"id","latitude","longitude","hotelname","countryid"
3,24.082284426992,32.887715399265,Sofitel Legend Old Cataract Aswan,EG
4,28.479180941035,34.49981957674,Tirana Dahab Resort,EG
5,27.106086,33.829103,Aladdin Beach Resort,EG
8,27.213800448298,33.841209806221,Aqua Fun,EG
9,27.240846416833,33.848536046974,Arabella Azur Resort,EG
10,27.241509667284,33.845314979553,Arabia Azur Resort,EG
11,27.251453,33.836651,Beirut,EG
12,27.102122033275,33.828242719174,SUNRISE Crystal Bay Resort - Grand Select,EG
13,27.080707657781,33.868165728836,Hurghada Coral Beach Hotel,EG
```

### Output

```
id searched,id data,meters,searched others,data others
10043289,554431,60,,map[countryid:GB hotelname:Olde Hall B & B]
10044206,1023038,300,,map[countryid:JP hotelname:Genius Resort]
10096121,517069,4,,map[countryid:FR hotelname:Villa Mauresque Arcachon]
10104857,559896,2,,map[countryid:CA hotelname:Brookside by Apex Accommodations]
10104962,1057206,51,,map[countryid:CA hotelname:Outlaws Inn by Apex Accommodations]
10217263,463869,1,,map[countryid:IT hotelname:Stella Di Rio]
10227010,218443,237,,map[countryid:CA hotelname:Stonegate Resort]
10260297,497838,219,,map[countryid:AR hotelname:Un Vistón Modern Mountain Home]
10296362,511525,32,,map[countryid:JP hotelname:Kyoto Fushimi Ohana]
10372374,502084,1,,map[countryid:NZ hotelname:The Church at Fox]
```
