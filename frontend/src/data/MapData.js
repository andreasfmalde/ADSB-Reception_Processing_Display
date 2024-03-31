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

export const initialView = {
    longitude: 10,
    latitude: 60.6,
    zoom: 5
};

