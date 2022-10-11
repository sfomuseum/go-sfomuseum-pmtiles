window.addEventListener('load', function(e){

    const map = L.map("map");
    map.setView([37.6143, -122.3828], 13);

    // note that it is 'sfo' and not 'sfo.pmtiles'
    // this is important
    
    var layer = protomaps.leafletLayer({url:"/sfo/{z}/{x}/{y}.mvt"});
    layer.addTo(map);
});
