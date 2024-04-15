window.addEventListener('load', function(e){

    // Remember: These get assigned at runtime in application/server/app.go

    var db_name = document.body.getAttribute("data-example-database");    
    var str_lat = document.body.getAttribute("data-example-latitude");
    var str_lon = document.body.getAttribute("data-example-longitude");
    var str_zoom = document.body.getAttribute("data-example-zoom");    

    var db_url = "/" + db_name + "/{z}/{x}/{y}.mvt";
    
    var lat = parseFloat(str_lat);
    var lon = parseFloat(str_lon);
    var zoom = parseInt(str_zoom);        
    
    const map = L.map("map");
    map.setView([lat, lon], zoom);

    var theme = 'light';
    var layer = protomapsL.leafletLayer({url: db_url, theme: theme});
    
    layer.addTo(map);
});
