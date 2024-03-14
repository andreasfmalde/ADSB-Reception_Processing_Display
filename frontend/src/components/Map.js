import React, { useRef, useEffect } from 'react';
import maplibregl from 'maplibre-gl';
import 'maplibre-gl/dist/maplibre-gl.css';
import './Map.css';


export default function Map() {
    const mapContainer = useRef(null);
    const map = useRef(null);



    useEffect(() => {
        if (map.current) return; // stops map from intializing more than once

        const style = {
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
    
        map.current = new maplibregl.Map({
          container: mapContainer.current,
          style: style,
          center: [10.68,60.79],
          zoom: 8
        });
      }, []);


      return (
        <div className="map-wrap">
            <div ref={mapContainer} className="map" />
        </div>
      );
}

