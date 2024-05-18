import { render, screen } from '@testing-library/react';
import  Sidebar from '../components/Sidebar';
import { testImage, geojson } from '../data/testData';
import '@testing-library/jest-dom/extend-expect';

describe('Sidebar component tests',()=>{
    
    test('the right message is displayed when aircraft props is null',()=>{
        render(<Sidebar aircraft={null} />);
        const message = screen.getByText('Select aircraft to view information');
        expect(message).toBeInTheDocument();
    });

    test('renders selected aircraft information', () => {
        
        

        render(<Sidebar aircraft={geojson.features[3]} image={testImage} />);

        const aircraftImage = screen.getByAltText('selected aircraft');
        expect(aircraftImage).toBeInTheDocument();

        const callsign = screen.getByText('SAS1812');
        expect(callsign).toBeInTheDocument();

        const icao = screen.getByText('45AC37');
        expect(icao).toBeInTheDocument();

        const photographer = screen.getByText('© Test Testesen');
        expect(photographer).toBeInTheDocument();

    });
    
});