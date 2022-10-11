window.addEventListener('load', function(e){

    var db_name = document.body.getAttribute("data-example-database");
    var db_url = "/" + db_name + "/{z}/{x}/{y}.mvt";
    
    var str_lat = document.body.getAttribute("data-example-latitude");
    var str_lon = document.body.getAttribute("data-example-longitude");
    var str_zoom = document.body.getAttribute("data-example-zoom");    

    var lat = parseFloat(str_lat);
    var lon = parseFloat(str_lon);
    var zoom = parseInt(str_zoom);        
    
    const map = L.map("map");
    map.setView([lat, lon], zoom);

    var layer = protomaps.leafletLayer({url: db_url});
    layer.addTo(map);
});
