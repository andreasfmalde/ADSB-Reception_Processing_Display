
import './App.css';
import { Navbar } from './components/Navbar';
import Map, {Marker} from 'react-map-gl/maplibre';
import 'maplibre-gl/dist/maplibre-gl.css';
import { useState } from 'react';
import { Sidebar } from './components/Sidebar';
import {style, geojson} from './data/MapData';
import { IoMdAirplane } from "react-icons/io";

function App() {
  const [viewport,setViewport] =  useState({
    longitude: 10,
    latitude: 60.6,
    zoom: 5
  });
  
  return (
    <div className="App">
      <Navbar />
      <div className="main-content">
        <Map
          initialViewState={viewport}
          minZoom={3}
          maxZoom={10}
          style={{width: 'calc(100vw - 300px)', height: 'calc(100vh - 78px)',gridColumn:'1/2'}}
          mapStyle={style}
          onMove={(e)=>{
            setViewport(e.viewState);
          }}
        >
          { geojson.features.map((p) =>(
          <div key={p.properties.icao}>
            <Marker 
              latitude={p.geometry.coordinates[0]}
              longitude={p.geometry.coordinates[1]}
              rotation={p.properties.track}
            >
            <IoMdAirplane 
              style={{
                color: '#c9a206',
                fontSize: '1.8em'
              }}
            />
          </Marker>
        </div>
        ))}
        </Map>
        <Sidebar />
      </div>
    </div>
  );
}

export default App;
