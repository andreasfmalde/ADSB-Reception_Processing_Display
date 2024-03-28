import { useState, useEffect, useRef } from 'react';
import Map, {Layer, Marker, Source} from 'react-map-gl/maplibre';
import { Sidebar } from './components/Sidebar';
import { Navbar } from './components/Navbar';
import {style, geojson, trail, trailLayer, initialView} from './data/MapData';
import { isInBounds, findAircraftByIcaoOrCallsign, trimAircraftList, callAPI } from './utils';
import { IoMdAirplane } from "react-icons/io";

import './App.css';
import 'maplibre-gl/dist/maplibre-gl.css';

function App() {
  const [viewport,setViewport] =  useState(initialView);
  const [aircraftJSON,setAircraftJSON] = useState(null);
  const [currentRender, setCurrentRender] = useState(null);
  const [selected, setSelected] = useState(null);
  const [selectedImg, setSelectedImg] = useState(null);
  const [historyTrail, setHistoryTrail] = useState(null);
  const map = useRef(null);


  const retrievePlanes = async () =>{
    try{
      const data = await callAPI(`${process.env.REACT_APP_SERVER}/aircraft/current/`);
      console.log(data)
      setAircraftJSON(data.features);
    }catch(error){
      console.log("Something went wrong")
    }
    
  }

  const retrieveImage = async (icao) =>{
    try{
      const data = await callAPI(`https://api.planespotters.net/pub/photos/hex/${icao}`);
      data.error ? setSelectedImg(null) : setSelectedImg(data.photos[0]);
    }catch(error){
      console.log("API retrieval failed")
    } 
  };

  const retrieveHistory = async (icao) =>{
    try{
      const data = await callAPI(`${process.env.REACT_APP_SERVER}/aircraft/history?icao=${icao}`);
      setHistoryTrail(data.features[0]);
    }catch(error){
      console.log("History not found")
    }
  }

  const retrieveSearch = (search) =>{
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

  const aircraftRenderFilter = (bounds) =>{
    let aircraftInBounds = aircraftJSON?.filter(p => isInBounds(p,bounds))
    setCurrentRender(trimAircraftList(aircraftInBounds));

  }

  useEffect(()=>{
    //  const seconds = 10;
    retrievePlanes();
    //setInterval(()=>retrievePlanes(),1000 * seconds );
  },[])
  
  return (
    <div className="App">
      <Navbar callback={retrieveSearch}/>
      <div className="main-content">
        <Map
          ref={map}
          initialViewState={viewport}
          minZoom={3}
          maxZoom={10}
          style={{width: 'calc(100vw - 300px)', height: 'calc(100vh - 78px)',gridColumn:'1/2'}}
          mapStyle={style}
          onLoad={(e)=>aircraftRenderFilter(e.target.getBounds())}
          onMoveEnd={(e)=>aircraftRenderFilter(e.target.getBounds())}
          onMove={(e)=>setViewport(e.viewState)}
          onClick={()=>{
            setHistoryTrail(null);
            setSelected(null);
          }}
        >
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
