// Make a call to an external API
export const callAPI = async (url) =>{
    const request = await fetch(url);
    return await request.json()
}