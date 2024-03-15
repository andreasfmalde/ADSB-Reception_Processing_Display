import logo from '../images/logo.png'
import './Navbar.css';

export const Navbar = () =>{


    return (
        <nav className="navigation">
            <img src={logo} alt="AirTrackr logo" className='logo' />
        </nav>
    );
}