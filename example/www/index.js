window.addEventListener('load', function(e){

    const map = L.map("map");
    map.setView([37.6143, -122.3828], 13);

    var db_name = document.body.getAttribute("data-example-database");
    var db_url = "/" + db_name + "/{z}/{x}/{y}.mvt";

    var layer = protomaps.leafletLayer({url: db_url});
    layer.addTo(map);
});
