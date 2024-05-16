import  Navbar  from "../components/Navbar";
import {render, screen, fireEvent} from "@testing-library/react";
import '@testing-library/jest-dom/extend-expect';

const mockFunction = jest.fn();

describe('Navbar component tests',()=>{
    test('main logo is rendered',()=>{
        render(<Navbar />);
        const image = screen.getByAltText("AirTrackr logo");
        expect(image).toBeInTheDocument();
    })


    test('searchbar is rendered',()=>{
        render(<Navbar />);
        const searchInput = screen.getByPlaceholderText("Search for callsign/icao...");
        const searchButton = screen.getByTestId("search-btn");
        expect(searchInput).toBeInTheDocument();
        expect(searchButton).toBeInTheDocument();
    })

    test('the search callback function is called when searched button is pressed',()=>{
        render(<Navbar callback={mockFunction}/>);
        const searchInput = screen.getByPlaceholderText("Search for callsign/icao...");
        const form = screen.getByTestId("form");

        fireEvent.change(searchInput, { target: { value: 'ABCDEF' } });
        fireEvent.submit(form);

        expect(mockFunction).toHaveBeenCalledTimes(1);
        expect(mockFunction).toHaveBeenCalledWith("ABCDEF");

    })
})