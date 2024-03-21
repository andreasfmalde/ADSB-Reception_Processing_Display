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