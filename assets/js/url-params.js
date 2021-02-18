function getUrlParameters() {
    const obj = {};
    const urlParams = new URLSearchParams(window.location.search);
    if (urlParams !== null) {
        const entries = urlParams.entries();
        for (pair of entries) {
            const key = pair[0];
            const value = pair[1];
            obj[key] = value;
        }
    }
    return obj;
}

function hasUrlParam(params, key) {
    return params.hasOwnProperty(key);
}

function getUrlParam(params, key, defaultValue) {
    if (params.hasOwnProperty(key) === true) {
        return params[key];
    }
    return defaultValue;
}

function getUrlParamAsBool(params, key, defaultValue) {
    if (params.hasOwnProperty(key) === true) {
        const value = params[key];
        return stringToBool(value, defaultValue);
    }
    return defaultValue;
}

function getUrlParamAsFloat(params, key, defaultValue) {
    if (params.hasOwnProperty(key) === true) {
        const value = params[key];
        return stringToFloat(value, defaultValue);
  }
    return defaultValue;
}

// see https://stackoverflow.com/questions/263965/how-can-i-convert-a-string-to-boolean-in-javascript
function stringToBool(str, defaultValue) {
    switch (str.toLowerCase().trim()) {
        case "1":
        case "true":
        case "yes":
        case "yay":
        case "yeah":
            return true;
        case "0":
        case "false":
        case "no":
        case "nay":
        case "nope":
        case null:
            return false;
        default:
            return defaultValue;
    }
}

function stringToFloat(str, defaultvalue) {
    const num = Number.parseFloat(str);
    return Number.isNaN(num) === false ? num : defaultValue;
}
