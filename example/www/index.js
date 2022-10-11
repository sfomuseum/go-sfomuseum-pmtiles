window.addEventListener('load', function(e){

    const map = L.map("map");

    map.setView([37.6143, -122.3828], 13);
    
    var layer = protomaps.leafletLayer({url:"/sfo.pmtiles/{z}/{x}/{y}.pbf"});
    layer.addTo(map);
});
