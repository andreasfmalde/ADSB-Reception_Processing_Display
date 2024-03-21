
import './App.css';
import { Navbar } from './components/Navbar';
import Map from 'react-map-gl/maplibre';
import 'maplibre-gl/dist/maplibre-gl.css';
import { useState } from 'react';
import { Sidebar } from './components/Sidebar';
import {style} from './data/MapData';

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
        </Map>
        <Sidebar />
      </div>
    </div>
  );
}

export default App;
