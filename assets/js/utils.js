// see https://stackoverflow.com/questions/49236100/copy-text-from-span-to-clipboard
function copyToClipboard(value) {
    const textArea = document.createElement('textarea');
    textArea.value = value;
    document.body.appendChild(textArea);
    textArea.select();
    textArea.setSelectionRange(0, 99999); // see https://www.w3schools.com/howto/howto_js_copy_clipboard.asp
    document.execCommand('copy');
    textArea.remove();
}

function convertLatLngToDMS(latitude, longitude, delimiter) {
    const lat = convertCoordinateToDMS(latitude);
    const lng = convertCoordinateToDMS(longitude);
    const degSymbol = '° ';
    return lat.degrees + degSymbol + lat.minutes + '\' ' + lat.seconds + '" ' + (latitude >= 0.0 ? 'N' : 'S')
        + delimiter + ' '
        + lng.degrees + degSymbol + lng.minutes + '\' ' + lng.seconds + '" ' + (longitude >= 0.0 ? 'E' : 'W');
}

// see https://stackoverflow.com/questions/37893131/how-to-convert-lat-long-from-decimal-degrees-to-dms-format
function convertCoordinateToDMS(coordinate) {
    const absolute = Math.abs(coordinate);
    const degrees = Math.floor(absolute);
    const minutesNotTruncated = (absolute - degrees) * 60;
    const minutes = Math.floor(minutesNotTruncated);
    const seconds = Number(((minutesNotTruncated - minutes) * 60.0).toFixed(1));
    return {degrees, minutes, seconds};
}

// see https://stackoverflow.com/questions/43167417/calculate-distance-between-two-points-in-leaflet
function calcDistance(from, to) {
    // returns distance in meters
    const degToRad = (deg) => deg * Math.PI / 180.0;
    const
        lat1 = degToRad(from[0]),
        lon1 = degToRad(from[1]),
        lat2 = degToRad(to[0]),
        lon2 = degToRad(to[1]);

    const dtLat = lat2 - lat1;
    const dtLon = lon2 - lon1;

    const a = Math.pow(Math.sin(dtLat * 0.5), 2) + Math.cos(lat1) * Math.cos(lat2) * Math.pow(Math.sin(dtLon * 0.5), 2);
    const c = 2 * Math.asin(Math.sqrt(a));
    const earthRadius = 6371.0 * 1000.0;
    return c * earthRadius;
}

// https://gist.github.com/demonixis/5188326
function toggleFullscreen(event) {
    let element = document.body;
    if (event instanceof HTMLElement) {
        element = event;
    }

    const isFullscreen = document.webkitIsFullScreen || document.mozFullScreen || false;
    element.requestFullScreen = element.requestFullScreen || element.webkitRequestFullScreen || element.mozRequestFullScreen || function () { return false; };
    document.cancelFullScreen = document.cancelFullScreen || document.webkitCancelFullScreen || document.mozCancelFullScreen || function () { return false; };
    isFullscreen ? document.cancelFullScreen() : element.requestFullScreen();
}

function getMapboxToken() {
    return 'pk.eyJ1IjoiZ3J1bXB5cGl4ZWwiLCJhIjoiY2trbWc5bWs4MmZhdTJ2bjdwbzYydmdyZCJ9.VRa2o7y--UXbK935iz7vvQ';
}

function getBingMapsKey() {
    return 'ApWBpzJqCe1z7pgOogc4fqVqlJldbTklgFFvzyaYCbmbuSEAQ_7mxOD_4ST3pA6i';
}

function getRandomGreeting() {
    const greetings = [
        'Ahlan',
        'At Your Service',
        'Ahoj',
        'Ciao',
        'Fly, you fools',
        'Hallo',
        'Hello',
        'Hai',
        'Hei',
        'Hej',
        'Hey',
        'Hi',
        'Hoi',
        'Hola',
        'Nǐ hǎo',
        'Olá',
        'Selam',
        'Salut',
        'Salve',
        'Yassou',
        'Yā, Yō',
    ];
    return greetings[Math.floor(Math.random() * greetings.length)] + '!';
}

function getRandomWatermark() {
    const watermarks = [
        'https://media.giphy.com/media/1YLcZOlQTKRmo/giphy.gif',
        'https://media.giphy.com/media/e5kbmb3wX3J1S/giphy.gif',
        'https://media.giphy.com/media/xRJZH4Ajr973y/giphy.gif',
        '/assets/img/github-mark-120px-plus.png',
        '/assets/img/github-mark-light-120px-plus.png',
    ];
    return watermarks[Math.floor(Math.random() * watermarks.length)];
}

function getRandomAirport() {
    const airports = [
        [65.653953, -18.073572], // BIAR
        [65.275289, -14.407159], // BIEG
        [63.981874, -22.585138], // BIKF
        [49.955209, -119.379941], // CYLW
        [44.236042, -78.356010], // CYPQ
        [43.631210, -79.392649], // CYTZ
        [45.458502, -73.754935], // CYUL
        [49.198174, -123.182909], // CYVR
        [51.130396, -114.015062], // CYYC
        [55.930905, -129.985748], // CZST
        [50.900214, 4.491274], // EBBR
        [50.047457, 8.574429], // EDDF
        [51.289167, 6.781183], // EDDL
        [50.870458, 7.126312], // EDDK
        [48.351439, 11.777140], // EDDM
        [60.313419, 24.963383], // EFHK
        [53.358943, -2.270278], // EGCC
        [51.397237, -3.337500], // EGFF
        [51.153031, -0.179936], // EGKK
        [51.468122, -0.449785], // EGLL
        [55.948236, -3.357371], // EGPH
        [52.299225, 4.753351], // EHAM
        [55.038504, -8.341528], // EIDL
        [53.421678, -6.236182], // EIDW
        [55.626197, 12.643986], // EKCH
        [49.632698, 6.215826], // ELLX
        [60.195993, 11.098937], // ENGM
        [59.643946, 17.927393], // ESSA
        [-33.966610, 18.598998], // FACT
        [-26.145187, 28.229308], // FAOR
        [30.108950, 31.400336], // HECA
        [0.044560, 32.441913], // HUEN
        [39.222338, -106.866782], // KASE
        // [33.637537, -84.439794], // KATL // This caused OSM to throw an exception...
        [39.857892, -104.672833], // KDEN
        [32.895095, -97.042235], // KDFW
        [59.349844, -151.927584], // KEB
        [40.635396, -73.781366], // KJFK
        [33.939690, -118.404956], // KLAX
        [28.435855, -81.313641], // KMCO
        [41.972954, -87.906999], // KORD
        [47.493300, -122.217610], // KRNT
        [32.734268, -117.201220], // KSAN
        [47.440480, -122.303866], // KSEA
        [34.849259, -111.791062], // KSEZ0
        [37.616258, -122.380741], // KSFO
        [37.953954, -107.901970], // KTEX
        [37.237503, -115.811926], // KXTA
        [41.419492, 19.715198], // LATI
        [41.289977, 2.079523], // LEBL
        [40.489425, -3.567624], // LEMD
        [39.547491, 2.736262], // LEPA
        [45.396843, 6.632939], // LFLJ
        [48.999596, 2.537087], // LFPG
        [37.928490, 23.939485], // LGAV
        [47.431108, 19.259556], // LHBP
        [40.902412, 9.518719], // LIEO
        [41.790567, 12.243056], // LIRF
        [45.615388, 8.722326], // LIMC
        [50.103618, 14.267351], // LKPR
        [31.999526, 34.876507], // LLBG
        [47.258157, 11.349303], // LOWI
        [48.116127, 16.570272], // LOWW
        [32.693853, -16.776296], // LPMA
        [38.771147, -9.131227], // LPPT
        [46.235005, 6.109182], // LSGG
        [36.897176, 30.798014], // LTAI
        [36.153008, -5.346761], // LXGB
        [14.061006, -87.218646], // MHTG
        [32.544717, -116.976464], // MMTJ
        [8.480211, -83.589702], // MRSN
        [-45.022708, 168.741582], // NZQN
        [33.822914, 35.491761], // OLBA
        [25.245723, 55.366214], // OMDB
        [35.544265, 139.785034], // RJTT
        [-22.815643, -43.247706], // SBGL
        [-0.127614, -78.358438], // SEQM
        [-22.819877, -43.243535], // SPGL
        [17.903753, -62.844720], // TFFJ
        [17.644494, -63.220311], // TNCS
        [55.977881, 37.418076], // UUEE
        [27.687352, 86.731662], // VNLK
        [13.687293, 100.750174], // VTBS
        [27.405689, 89.421487], // VQPR
        [-6.131338, 106.657723], // WIII
        [-3.623944, 136.594537], // WX53
        [-33.934871, 151.177930], // YSSY
    ];
    if (airports.length > 0) {
        return airports[Math.floor(Math.random() * airports.length)];
    } else {
        return [0, 0];
    }
}