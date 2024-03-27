import './components.css'

export const Sidebar = (props) =>{
    return(
        <div className='Sidebar'>
        {props?.aircraft == null ? 
            <div className='sidebar-unselected'>
                Select aircraft to view information
            </div>:
            // Render aircraft
            <div className='aircraft-info'>
                <img src={props?.image != null ? props?.image.thumbnail_large.src : "https://eagle-sensors.com/wp-content/uploads/unavailable-image.jpg" } 
                alt='selected aircraft'/>
                {props?.image != null ? <a href={props?.image.link}><span>&copy; {props?.image.photographer}</span></a> : ""}
                <div className='callsign'>{props?.aircraft.properties.callsign}</div>
                <div className='property'><div className='label'>ALTITUDE:</div> {props?.aircraft.properties.altitude} feet</div>
                <div className='property'><div className='label'>SPEED:</div> {props?.aircraft.properties.speed} knots</div>
                <div className='property'><div className='label'>VERTICAL SPEED:</div> {props?.aircraft.properties.vspeed} feet/min</div>
                <div className='property'><div className='label'>TRACK:</div> {props?.aircraft.properties.track} &deg;</div>
                <div className='property'><div className='label'>POSITION:</div> Lat: {props?.aircraft.geometry.coordinates[0]} N
                <br/> Long: {props?.aircraft.geometry.coordinates[1]} W </div>
                <div className='property'><div className='label'>ICAO:</div> {props?.aircraft.properties.icao}</div>
                <div className='property'><div className='label'>TIMESTAMP:</div> {props?.aircraft.properties.timestamp}</div>
            </div>
        }
        </div>
    );
}