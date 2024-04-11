import './components.css';
import { IoMdSearch, IoMdSettings } from "react-icons/io";
import logoWhite from '../assets/logo_white.png';
import { useState, useRef, useEffect } from 'react';

export const Navbar = (props) =>{

    const [open, setOpen] = useState(false);
    const [historyHour, setHistoryHour] = useState(1);
    const [currentCustom, setCurrentCustom] = useState(true);
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
                <button onClick={()=>setOpen(!open)} ref={dropDownButton}>
                    <IoMdSettings className='settings-icon' />
                </button>
                <div className={`drop-down-window ${open ? 'active':'inactive'}`} ref={dropDownRef}>
                    <h3>HISTORY TRAILS:</h3>
                    <p className='label'>Select how many hours in the past the history trails will show. <br /> Choose to set a custom amount, or show all history data.</p>
                    <hr />
                    <input type='radio' id='custom-radio' checked={currentCustom}
                    onClick={()=>{
                        setCurrentCustom(true)
                        props.trail(`${historyHour}`)
                    }}/>
                    <label htmlFor="custom-radio">Custom hours</label>

                    <div className='custom-history' >
                        <input type='button' value='-'
                            onClick={()=>{
                                if(historyHour > 1){
                                    setHistoryHour(historyHour-1)
                                    props.trail(`${historyHour-1}`)
                                    console.log(`${historyHour-1}`)
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
                                    console.log(`${historyHour+1}`)
                                }
                            }}
                            disabled={!currentCustom}
                        />
                    </div>
                    
                    <input type='radio' id='option-4' value='all'
                    checked={!currentCustom}
                    onClick={e=>{
                        setCurrentCustom(false)
                        props.trail(e.target.value)
                        }} 
                        />
                    <label htmlFor="option-4">Show all</label>
                </div>
            </div>
        </nav>
    );
}