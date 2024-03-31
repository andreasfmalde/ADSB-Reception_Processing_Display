import { callAPI } from "./apiUtils";

describe('callAPI tests',()=>{

    test('the right url is fetched and a valid json is retrieved',async() =>{
        global.fetch = jest.fn();
        jest.spyOn(global,'fetch').mockImplementationOnce((url)=>{
            const output = {
                msg: "This is json",
                url: url,
            }
            return Promise.resolve({
                ok:true,
                json: () => Promise.resolve(output),
            })
        });
        try{
            const res = await callAPI('http://someaddress.no');
            expect(res.msg).toStrictEqual("This is json");
            expect(res.url).toStrictEqual("http://someaddress.no")
        }catch(error){
            fail('something went wrong');
        }
    });
});