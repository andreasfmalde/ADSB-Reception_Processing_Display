
import './App.css';
import { Navbar } from './components/Navbar';
import Map from 'react-map-gl/maplibre';
import 'maplibre-gl/dist/maplibre-gl.css';
import { useState } from 'react';

function App() {
  const [viewport,setViewport] =  useState({
    longitude: 10,
    latitude: 60.6,
    zoom: 5
  });
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
  return (
    <div className="App">
      <Navbar />
      <div className="main-content">
        <Map
          initialViewState={viewport}
          minZoom={3}
          maxZoom={10}
          style={{width: '1000px', height: '600px'}}
          mapStyle={style}
          onMove={(e)=>{
            setViewport(e.viewState);
          }}
        >
        </Map>

      </div>
    </div>
  );
}

export default App;
