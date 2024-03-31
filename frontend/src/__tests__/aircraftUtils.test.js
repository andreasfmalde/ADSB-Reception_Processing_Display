import { findAircraftByIcaoOrCallsign, isInBounds, trimAircraftList} from "../utils/aircraft/aircraftUtils";
import { geojson, mapBounds } from "../data/TestData";



describe('findAircraftByIcaoOrCallsign tests',()=>{

   test('returns null when aircraft list is null', () => {
        const result = findAircraftByIcaoOrCallsign("A3F9F",null);
        expect(result).toBe(null); 
    }); 

    test('returns aircraft when search term matches callsign',()=>{
        const result = findAircraftByIcaoOrCallsign("SAS1812", geojson.features);
        expect(result).toStrictEqual(geojson.features[3]);
    });

    test('returns aircraft when search term matches icao',()=>{
        const result = findAircraftByIcaoOrCallsign("45AC32", geojson.features);
        expect(result).toStrictEqual(geojson.features[2]);
    });

    test('returns null when search term does not match any aircrafts in the list',()=>{
        const result = findAircraftByIcaoOrCallsign("NOT123", geojson.features);
        expect(result).toBe(null);
    });

    test('returns null when aircraft list is empty',()=>{
        const result = findAircraftByIcaoOrCallsign("45AC32", []);
        expect(result).toBe(null);
    });

});


describe('isInBounds tests',()=>{
    let aircraftClone;

    beforeEach(()=>{
        aircraftClone = {
            type: 'Feature',
            properties: {
              icao: '45AC37',
              callsign: 'SAS1812',
              altitude: 37000,
              speed: 494,
              track: 39,
              vspeed: 0,
              timestamp: '2024-03-15T12:49:19Z'
            },
            geometry: { coordinates: [ 57.82393, 16.3162 ], type: 'Point' }
        };
    })

    test('return true when aircraft is within map bounds',()=>{
        const result = isInBounds(geojson.features[3],mapBounds);
        expect(result).toStrictEqual(true);
    });
    test('return false when aircraft is not within map bounds',()=>{
        const result = isInBounds(geojson.features[0],mapBounds);
        expect(result).toStrictEqual(false);
    });
    test('return false when aircraft is null or undefined',()=>{
        let result = isInBounds(null,mapBounds);
        expect(result).toStrictEqual(false);
        result = isInBounds(undefined,mapBounds);
        expect(result).toStrictEqual(false);
    });
    test('return false when map bounds is null or undefined',()=>{
        let result = isInBounds(geojson.features[3],null);
        expect(result).toStrictEqual(false);
        result = isInBounds(geojson.features[3],undefined);
        expect(result).toStrictEqual(false);
    });
    test('return false when aircraft is within long bounds, but not lat bounds',()=>{
        const result = isInBounds(aircraftClone,mapBounds);
        expect(result).toStrictEqual(true);
        // Aircraft is too north in relation to the map bounds
        aircraftClone.geometry.coordinates[0] = 65;
        const resultTooNorth = isInBounds(aircraftClone,mapBounds);
        expect(resultTooNorth).toStrictEqual(false);
        // Aircraft is too south in relation to the map bounds
        aircraftClone.geometry.coordinates[0] = 50;
        const resultTooSouth = isInBounds(aircraftClone,mapBounds);
        expect(resultTooSouth).toStrictEqual(false);

    });

    test('return false when aircraft is within lat bounds, but not long bounds',()=>{
        const result = isInBounds(aircraftClone,mapBounds);
        expect(result).toStrictEqual(true);
        // Aircraft is too west in relation to the map bounds
        aircraftClone.geometry.coordinates[1] = -10;
        const resultTooWest = isInBounds(aircraftClone,mapBounds);
        expect(resultTooWest).toStrictEqual(false);
        // Aircraft is too east in relation to the map bounds
        aircraftClone.geometry.coordinates[0] = 25;
        const resultTooEast = isInBounds(aircraftClone,mapBounds);
        expect(resultTooEast).toStrictEqual(false);
    });   
});

describe('trimAircraftList tests', ()=>{
    test('returns null when aircraft list is null or undefined',()=>{
        let result = trimAircraftList(null);
        expect(result).toBe(null);
        result = trimAircraftList(undefined);
        expect(result).toBe(null);
    });

    test('trims the length of the list to be within proper limits',()=>{

        let mockAirCraftList;
        let result;

        const testValues = [5000,3000,2200,1500,1000,700]
        testValues.forEach(val =>{
            mockAirCraftList = new Array(val).fill('ac');
            result = trimAircraftList(mockAirCraftList).length;
            // The length of the list should be reduced.
            // The new length should not be above 900
            expect(result).toBeLessThan(val);
            expect(result).toBeLessThan(900);
        })
        // When under 500 aircrafts in the list, no trimming
        // is done and all aircrafts are returned
        mockAirCraftList = new Array(499).fill('ac');
        result = trimAircraftList(mockAirCraftList).length;
        expect(result).toStrictEqual(499)

    })
});