import './components.css';
import { IoMdSearch } from "react-icons/io";
import logoWhite from '../assets/logo_white.png';

export const Navbar = (props) =>{


    return (
        <nav className="navigation">
            <img src={logoWhite} alt="AirTrackr logo" className='logo' />
            <form className='search-field'
                onSubmit={e =>{
                    e.preventDefault();
                    props.callback(e.target.searchbar.value);
                    console.log(e.target.searchbar.value)
                }}
                >
                <input type="text" placeholder='Search for callsign/icao...' name="searchbar"/>
                <button type='submit' data-testid="search-btn">
                    <IoMdSearch />
                </button>
            </form>
        </nav>
    );
}