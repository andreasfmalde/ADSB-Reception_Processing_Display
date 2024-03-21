
import './App.css';
import { Navbar } from './components/Navbar';
import Map, {Marker} from 'react-map-gl/maplibre';
import 'maplibre-gl/dist/maplibre-gl.css';
import { useState, useEffect } from 'react';
import { Sidebar } from './components/Sidebar';
import {style, geojson} from './data/MapData';
import { IoMdAirplane } from "react-icons/io";

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
      const response = await fetch("http://localhost:8080/aircraft/current/");// http://129.241.150.147:8080/aircraft/current/
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
      if (data.error){
        console.log("AYA")
        setSelectedImg(null)
      }else{
        setSelectedImg(data.photos[0])
      }
    }catch(error){
      console.log("API retrieval failed")
    } 
  };

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
        </Map>
        <Sidebar aircraft={selected} image={selectedImg} />
      </div>
    </div>
  );
}

export default App;
