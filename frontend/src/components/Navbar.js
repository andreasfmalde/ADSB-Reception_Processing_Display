import './components.css';
import { IoMdSearch, IoMdSettings, IoIosCloseCircleOutline  } from "react-icons/io";
import logoWhite from '../assets/logo_white.png';
import { useState, useRef, useEffect } from 'react';

export const Navbar = (props) =>{

    const [open, setOpen] = useState(false);
    const [historyHour, setHistoryHour] = useState(1);
    const [currentCustom, setCurrentCustom] = useState(true);
    const [searchOpen, setSearchOpen] = useState(false);
    const dropDownRef = useRef();
    const dropDownButton = useRef();

    useEffect(()=>{

        document.addEventListener('mousedown',e =>{
            if(open && !dropDownRef.current.contains(e.target) 
            && !dropDownButton.current.contains(e.target)){
                setOpen(false)
            }
        });
    })
    return (
        <nav className="navigation">
            {/* LOGO */}
            <img src={logoWhite} alt="AirTrackr logo" className='logo' />
            {/* SEARCH BAR */}
            <form className={`search-field ${searchOpen ? 'search-active':''}`}
                data-testid='form'
                onSubmit={e =>{
                    e.preventDefault();
                    props.callback(e.target.querySelector('[name="searchbar"]').value);
                    e.target.querySelector('[name="searchbar"]').value = "";
                }}
                >
                <input type="text" placeholder='Search for callsign/icao...' name="searchbar"/>
                <button type='submit' data-testid='search-btn'>
                    <IoMdSearch />
                </button>
            </form>
            {/* SEARCH BAR BUTTON FOR SMALL/MEDIUM SCREENS */}
            <button className='search-btn-medium-small'
                onClick={()=>{
                    setSearchOpen(!searchOpen)
                    setOpen(false)
                }}
            >
                {searchOpen ? <IoIosCloseCircleOutline /> :<IoMdSearch />}
            </button>
            {/* HISTORY TRAIL DROP DOWN MENU */}
            <div className='drop-down'>
                <button 
                    onClick={()=>{
                        setOpen(!open)
                        setSearchOpen(false)
                    }} 
                    ref={dropDownButton}
                >
                    <IoMdSettings />
                </button>
                {/* DROP DOWN WINDOW */}
                <div className={`drop-down-window ${open ? 'active':'inactive'}`} ref={dropDownRef}>
                    <h3>HISTORY TRAILS:</h3>
                    <p className='label'>Select how many hours in the past the history trails will show. <br />
                    Choose to set a custom amount, or show all history data.</p>
                    <hr />
                    {/* CUSTOM SECTION */}
                    <input type='radio' id='custom-radio' checked={currentCustom}
                    onChange={() => setCurrentCustom(true)} 
                    onClick={()=>{
                        props.trail(`${historyHour}`)
                    }}/>
                    <label htmlFor="custom-radio">Custom hours</label>

                    <div className='custom-history' >
                        <input type='button' value='-'
                            onClick={()=>{
                                if(historyHour > 1){
                                    setHistoryHour(historyHour-1)
                                    props.trail(`${historyHour-1}`)
                                }
                            }}
                            disabled={!currentCustom}
                        />
                        <input type="text" name="custom-nr"
                         readOnly value={historyHour}
                         disabled={!currentCustom}/>
                        <input type='button' value='+'
                            onClick={()=>{
                                if(historyHour < 24){
                                    setHistoryHour(historyHour+1)
                                    props.trail(`${historyHour+1}`)
                                }
                            }}
                            disabled={!currentCustom}
                        />
                    </div>
                    {/* ALL DATA SECTION */}
                    <input type='radio' id='option-4' value='all'
                    checked={!currentCustom}
                    onChange={() => setCurrentCustom(false)} 
                    onClick={e=>{
                        props.trail(e.target.value)
                        }} 
                        />
                    <label htmlFor="option-4">Show all</label>
                </div>
            </div>
        </nav>
    );
}