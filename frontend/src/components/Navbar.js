import './components.css';
import { IoMdSearch, IoMdSettings } from "react-icons/io";
import logoWhite from '../assets/logo_white.png';
import { useState } from 'react';

export const Navbar = (props) =>{

    const [open, setOpen] = useState(false);

    return (
        <nav className="navigation">
            <img src={logoWhite} alt="AirTrackr logo" className='logo' />
            <form className='search-field'
                data-testid='form'
                onSubmit={e =>{
                    e.preventDefault();
                    props.callback(e.target.querySelector('[name="searchbar"]').value);
                }}
                >
                <input type="text" placeholder='Search for callsign/icao...' name="searchbar"/>
                <button type='submit' data-testid='search-btn'>
                    <IoMdSearch />
                </button>
            </form>
            <div className='drop-down'>
                <button onClick={()=>setOpen(!open)}>
                    <IoMdSettings className='settings-icon' />
                </button>
                <div className={`drop-down-window ${open ? 'active':'inactive'}`}>
                    <h3>Set history trail length:</h3>
                    <hr />
                    <form 
                        onSubmit={e=>{
                            e.preventDefault();
                            props.trail(e.target.trail.value)
                        }}
                    >
                        <input type='radio' id='option-1' name="trail" value='1' />
                        <label htmlFor="option-1">1 Hour</label>
                        <br />
                        <input type='radio' id='option-2' name="trail" value='5' />
                        <label htmlFor="option-2">5 Hours</label>
                        <br />
                        <input type='radio' id='option-3' name="trail" value='24' />
                        <label htmlFor="option-3">1 Day</label>
                        <br />
                        <input type='radio' id='option-4' name="trail" value='all' />
                        <label htmlFor="option-4">All</label>
                        <br />
                        <input type="submit" value='Save'/>

                    </form>

                </div>
            </div>
        </nav>
    );
}