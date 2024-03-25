import { useState, useEffect } from 'react';
import Map, {Layer, Marker, Source} from 'react-map-gl/maplibre';
import { Sidebar } from './components/Sidebar';
import { Navbar } from './components/Navbar';
import {style, geojson, trail, trailLayer} from './data/MapData';
import { IoMdAirplane } from "react-icons/io";

import './App.css';
import 'maplibre-gl/dist/maplibre-gl.css';

function App() {
  const [viewport,setViewport] =  useState({
    longitude: 10,
    latitude: 60.6,
    zoom: 5
  });
  const [aircraftJSON,setAircraftJSON] = useState(null);
  const [currentRender, setCurrentRender] = useState(null);
  const [selected, setSelected] = useState(null);
  const [selectedImg, setSelectedImg] = useState(null);
  const [historyTrail, setHistoryTrail] = useState(null);

  const isInBounds = (p,mapBounds) =>{
      
    if (p.geometry.coordinates[0] > mapBounds._ne.lat || p.geometry.coordinates[0] < mapBounds._sw.lat ){
      return false;
    }
    if (p.geometry.coordinates[1] > mapBounds._ne.long || p.geometry.coordinates[1] < mapBounds._sw.long ){
      return false;
    }
    return true
  }

  const retrievePlanes = async () =>{
    try{
      const response = await fetch("http://localhost:8080/aircraft/current/"); // http://129.241.150.147:8080/aircraft/current/
      const data = await response.json()
      setAircraftJSON(data.features);
    }catch(error){
      console.log("Something went wrong")
    }
    
  }

  const retrieveImage = async (icao) =>{
    try{
      const response = await fetch(`https://api.planespotters.net/pub/photos/hex/${icao}`);
      const data = await response.json();
      data.error ? setSelectedImg(null) : setSelectedImg(data.photos[0]);
    }catch(error){
      console.log("API retrieval failed")
    } 
  };

  const retrieveHistory = async (icao) =>{
    try{
      const response = await fetch(`http://localhost:8080/aircraft/history?icao=${icao}`);
      const data = await response.json();
      setHistoryTrail(data.features[0]);
    }catch(error){
      console.log("History not found")
    }
  }

  useEffect(()=>{
    //  const seconds = 10;
    retrievePlanes();
    //setInterval(()=>retrievePlanes(),1000 * seconds );
  },[])
  
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
          onMoveEnd={(e)=>{
            let aircraftInBounds = aircraftJSON?.filter(p => isInBounds(p,e.target.getBounds()))
            if(aircraftInBounds !== undefined){
              if (aircraftInBounds.length > 3500){
                aircraftInBounds = aircraftInBounds.filter(() => Math.random() > 0.9)
              }else if (aircraftInBounds.length > 2500){
                aircraftInBounds = aircraftInBounds.filter(() => Math.random() > 0.8)
              }else if (aircraftInBounds.length > 2000){
                aircraftInBounds = aircraftInBounds.filter(() => Math.random() > 0.7)
              }else if (aircraftInBounds.length > 1250){
                aircraftInBounds = aircraftInBounds.filter(() => Math.random() > 0.5)
              }else if (aircraftInBounds.length > 900){
                aircraftInBounds = aircraftInBounds.filter(() => Math.random() > 0.3)
              }else if (aircraftInBounds.length > 500){
                aircraftInBounds = aircraftInBounds.filter(() => Math.random() > 0.15)
              }
            }
            
            setCurrentRender(aircraftInBounds);
          }}
        >
          {currentRender?.map((p) =>(
          <div key={p.properties.icao}>
            <Marker 
              latitude={p.geometry.coordinates[0]}
              longitude={p.geometry.coordinates[1]}
              rotation={p.properties.track}
              onClick={()=>{
                if(selected?.properties.icao !== p.properties.icao){
                  setSelected(p);
                  retrieveImage(p.properties.icao);
                  retrieveHistory(p.properties.icao);
                }
                
              }}
            >
            <IoMdAirplane 
              style={{
                color: selected?.properties.icao===p.properties.icao ? '#b50404' : '#c9a206',
                fontSize: '1.8em'
              }}
            />
          </Marker>
        </div>
        ))}
        {selected === null ? "":
          <Source id='trail' type='geojson' data={historyTrail}>
            <Layer {... trailLayer} />
          </Source>
        }
        </Map>
        <Sidebar aircraft={selected} image={selectedImg} />
      </div>
    </div>
  );
}

export default App;
