import { useState, useEffect, useRef } from 'react';
import Map, {Layer, Marker, Source} from 'react-map-gl/maplibre';
import { Sidebar } from './components/Sidebar';
import { Navbar } from './components/Navbar';
import {style, trailLayer, initialView} from './data/MapData';
import { isInBounds, findAircraftByIcaoOrCallsign, trimAircraftList, callAPI } from './utils';
import { IoMdAirplane } from "react-icons/io";

import './App.css';
import 'maplibre-gl/dist/maplibre-gl.css';

// Main point of the application 
function App() {
  const [viewport,setViewport] =  useState(initialView);
  const [aircraftJSON,setAircraftJSON] = useState(null);
  const [currentRender, setCurrentRender] = useState(null);
  const [selected, setSelected] = useState(null);
  const [selectedImg, setSelectedImg] = useState(null);
  const [historyTrail, setHistoryTrail] = useState(null);
  const [currentBounds, setCurrentBounds] = useState(null);
  const map = useRef(null);

  // Retrieve aircrafts from API and update the current aircraft list
  const retrievePlanes = async () =>{
    try{
      const data = await callAPI(`${process.env.REACT_APP_SERVER}/aircraft/current/`);
      setAircraftJSON(data.features);
    }catch(error){
      console.log("No aircrafts are fetched")
    }
    
  }
  // Retrieve an aircraft image from third party API based on
  // aircraft icao
  const retrieveImage = async (icao) =>{
    try{
      const data = await callAPI(`https://api.planespotters.net/pub/photos/hex/${icao}`);
      data.error ? setSelectedImg(null) : setSelectedImg(data.photos[0]);
    }catch(error){
      console.log("API retrieval failed")
    } 
  };

  // Retrieve historical location coordinates for an aircraft
  // with the specified icao
  const retrieveHistory = async (icao) =>{
    try{
      const data = await callAPI(`${process.env.REACT_APP_SERVER}/aircraft/history/${icao}`);
      setHistoryTrail(data.features[0]);
    }catch(error){
      console.log("History not found")
    }
  }

  // Search for aircraft based on icao or callsign and 
  // make the map fly to the aircrafts' location
  const searchForAircraft = (search) =>{
    if (search === null || search === undefined || search.length < 3 || search.length > 15){
      return
    }
    let aircraft = findAircraftByIcaoOrCallsign(search,aircraftJSON);
    if (aircraft !== null){
      setSelected(aircraft);
      retrieveImage(aircraft.properties.icao);
      retrieveHistory(aircraft.properties.icao);
      map.current.flyTo({center:[aircraft.geometry.coordinates[1],aircraft.geometry.coordinates[0]],zoom:9})
    }
    
  }
  
  // Update/render the aircrafts on the map every time updated
  // aircraft information is retrieved from external API or when
  // the map is moved to a new location
  useEffect(()=>{
    if(currentBounds !== null && aircraftJSON !== null){
      let aircraftInBounds = aircraftJSON?.filter(p => isInBounds(p,currentBounds));
      setCurrentRender(trimAircraftList(aircraftInBounds));
    }
    
  },[aircraftJSON, currentBounds])
  
  return (
    <div className="App">
      <Navbar callback={searchForAircraft}/>
      <div className="main-content">
        <Map
          className='main-map'
          ref={map}
          initialViewState={viewport}
          minZoom={3}
          maxZoom={10}
          mapStyle={style}
          onLoad={(e)=>{
            retrievePlanes();
            window.setInterval(()=>retrievePlanes(),30000)
            setCurrentBounds(e.target.getBounds())
          }}
          onMoveEnd={(e)=>setCurrentBounds(e.target.getBounds())}
          onMove={(e)=>setViewport(e.viewState)}
          onClick={()=>{
            setHistoryTrail(null);
            setSelected(null);
          }}
        >
          {/* Render the aircrafts currently within the viewport */}
          {currentRender?.map((p) =>(
          <div key={p.properties.icao}>
            <Marker 
              latitude={p.geometry.coordinates[0]}
              longitude={p.geometry.coordinates[1]}
              rotation={p.properties.track}
              onClick={e=>{
                e.originalEvent.stopPropagation();
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
        {/* Render the history trail behind a selected aircraft */}
        {selected === null || historyTrail === null ? "":
          <Source id='trail' type='geojson' data={historyTrail}>
            <Layer {... trailLayer} />
          </Source>
        }
        </Map>
        <Sidebar aircraft={selected} image={selectedImg}/>
      </div>
    </div>
  );
}

export default App;
