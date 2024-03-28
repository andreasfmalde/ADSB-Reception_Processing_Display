export const style = {
    "version": 8,
    "sources": {
      "osm": {
        "type": "raster",
        "tiles": ["https://tile.openstreetmap.org/{z}/{x}/{y}.png"],
        "tileSize": 256,
        "attribution": "&copy; OpenStreetMap Contributors",
        "maxzoom": 19
      }
    },
    "layers": [
      {
        "id": "osm",
        "type": "raster",
        "source": "osm" // This must match the source key above
      }
    ]
};

export const initialView = {
    longitude: 10,
    latitude: 60.6,
    zoom: 5
};


export const geojson = {
    "type": "FeatureCollection",
    "features": [
      {
        "type": "Feature",
        "properties": {
          "icao": "8BA5",
          "callsign": "OPM007",
          "altitude": 6175,
          "speed": 231,
          "track": 256,
          "vspeed": 256,
          "timestamp": "2024-03-15T12:49:19Z"
        },
        "geometry": {
          "coordinates": [
            40.785812,
            -74.27357
          ],
          "type": "Point"
        }
      },
      {
        "type": "Feature",
        "properties": {
          "icao": "8DF8",
          "callsign": "SFR345",
          "altitude": 6100,
          "speed": 239,
          "track": 324,
          "vspeed": 2496,
          "timestamp": "2024-03-15T12:49:18Z"
        },
        "geometry": {
          "coordinates": [
            -29.51123,
            31.118986
          ],
          "type": "Point"
        }
      },
      {
        "type": "Feature",
        "properties": {
          "icao": "45AC32",
          "callsign": "SAS4632",
          "altitude": 33975,
          "speed": 509,
          "track": 6,
          "vspeed": 0,
          "timestamp": "2024-03-15T12:49:18Z"
        },
        "geometry": {
          "coordinates": [
            54.06299,
            10.196411
          ],
          "type": "Point"
        }
      },
      {
        "type": "Feature",
        "properties": {
          "icao": "45AC37",
          "callsign": "SAS1812",
          "altitude": 37000,
          "speed": 494,
          "track": 39,
          "vspeed": 0,
          "timestamp": "2024-03-15T12:49:19Z"
        },
        "geometry": {
          "coordinates": [
            57.82393,
            16.3162
          ],
          "type": "Point"
        }
      }
    ]
};

export const trail = {
	"type": "FeatureCollection",
	"features": [
		{
			"type": "Feature",
			"properties": {
				"icao": "485086"
			},
			"geometry": {
				"coordinates": [
					[
						8.651145,
						61.244797
					],
					[
						8.66017,
						61.263016
					],
					[
						8.668998,
						61.28087
					],
					[
						8.677139,
						61.297256
					],
					[
						8.686458,
						61.3163
					],
					[
						8.693422,
						61.33026
					],
					[
						8.699504,
						61.342438
					],
					[
						8.718828,
						61.38121
					],
					[
						8.727461,
						61.398605
					],
					[
						8.746589,
						61.436966
					],
					[
						8.755613,
						61.455
					],
					[
						8.774153,
						61.49208
					],
					[
						8.782,
						61.50769
					],
					[
						8.789161,
						61.52188
					],
					[
						8.797792,
						61.539047
					],
					[
						8.817705,
						61.578598
					],
					[
						8.82516,
						61.593475
					],
					[
						8.833989,
						61.61087
					],
					[
						8.857924,
						61.65413
					],
					[
						8.890686,
						61.698715
					],
					[
						8.899514,
						61.70906
					],
					[
						8.911776,
						61.72293
					],
					[
						8.925705,
						61.73845
					],
					[
						8.954249,
						61.770355
					],
					[
						9.048615,
						61.87509
					]
				],
				"type": "LineString"
			}
		}
	]
};

export const trailLayer = {
  'id': 'trailLayer',
  'type': 'line',
  'source': 'trail',
  'layout': {
      'line-join': 'round',
      'line-cap': 'round'
  },
  'paint': {
      'line-color': '#000',
      'line-width': 3
  }
};