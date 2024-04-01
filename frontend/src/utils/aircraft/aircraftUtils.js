// Determine if the location of an aircraft is inside the 
// the map bounds
export const isInBounds = (ac,mapBounds) =>{
  if (ac === null || ac === undefined || mapBounds === null || mapBounds === undefined){
    return false;
  }
  if (ac.geometry.coordinates[0] > mapBounds._ne.lat || ac.geometry.coordinates[0] < mapBounds._sw.lat ){
    return false;
  }
  if (ac.geometry.coordinates[1] > mapBounds._ne.lng || ac.geometry.coordinates[1] < mapBounds._sw.lng ){
    return false;
  }
  return true
}

// Find and return an aircraft with a given icao or callsign from 
// a list, else return null
export const findAircraftByIcaoOrCallsign = (search, aircrafts) =>{
  if (aircrafts !== null){
    for (let ac of aircrafts){
      if (ac.properties.icao === search || ac.properties.callsign === search){
        return ac
      }
    }
  }
  return null;
}

// Shorten the number of aircrafts based on the current amount of
// aircrafts in the list
export const trimAircraftList = (aircrafts) =>{
  if(aircrafts !== undefined && aircrafts !== null){
      if (aircrafts.length > 3500){
        aircrafts = aircrafts.filter(() => Math.random() > 0.9)
      }else if (aircrafts.length > 2500){
        aircrafts = aircrafts.filter(() => Math.random() > 0.8)
      }else if (aircrafts.length > 2000){
        aircrafts = aircrafts.filter(() => Math.random() > 0.7)
      }else if (aircrafts.length > 1250){
        aircrafts = aircrafts.filter(() => Math.random() > 0.5)
      }else if (aircrafts.length > 900){
        aircrafts = aircrafts.filter(() => Math.random() > 0.3)
      }else if (aircrafts.length > 500){
        aircrafts = aircrafts.filter(() => Math.random() > 0.15)
      }
      return aircrafts;
  }
  return null;
}