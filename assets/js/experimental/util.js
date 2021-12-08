class Util {
    static randomNumber(min, max) {
        return Math.random() * (max - min) + min;
    }
    static randomInt(min, max) {
        return Math.floor(this.randomNumber(min, max));
    }
    static isNumber(x) {
        return typeof x === "number";
    }
}