/******/ (function(modules) { // webpackBootstrap
/******/ 	// The module cache
/******/ 	var installedModules = {};
/******/
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/
/******/ 		// Check if module is in cache
/******/ 		if(installedModules[moduleId]) {
/******/ 			return installedModules[moduleId].exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = installedModules[moduleId] = {
/******/ 			i: moduleId,
/******/ 			l: false,
/******/ 			exports: {}
/******/ 		};
/******/
/******/ 		// Execute the module function
/******/ 		modules[moduleId].call(module.exports, module, module.exports, __webpack_require__);
/******/
/******/ 		// Flag the module as loaded
/******/ 		module.l = true;
/******/
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/
/******/
/******/ 	// expose the modules object (__webpack_modules__)
/******/ 	__webpack_require__.m = modules;
/******/
/******/ 	// expose the module cache
/******/ 	__webpack_require__.c = installedModules;
/******/
/******/ 	// define getter function for harmony exports
/******/ 	__webpack_require__.d = function(exports, name, getter) {
/******/ 		if(!__webpack_require__.o(exports, name)) {
/******/ 			Object.defineProperty(exports, name, { enumerable: true, get: getter });
/******/ 		}
/******/ 	};
/******/
/******/ 	// define __esModule on exports
/******/ 	__webpack_require__.r = function(exports) {
/******/ 		if(typeof Symbol !== 'undefined' && Symbol.toStringTag) {
/******/ 			Object.defineProperty(exports, Symbol.toStringTag, { value: 'Module' });
/******/ 		}
/******/ 		Object.defineProperty(exports, '__esModule', { value: true });
/******/ 	};
/******/
/******/ 	// create a fake namespace object
/******/ 	// mode & 1: value is a module id, require it
/******/ 	// mode & 2: merge all properties of value into the ns
/******/ 	// mode & 4: return value when already ns object
/******/ 	// mode & 8|1: behave like require
/******/ 	__webpack_require__.t = function(value, mode) {
/******/ 		if(mode & 1) value = __webpack_require__(value);
/******/ 		if(mode & 8) return value;
/******/ 		if((mode & 4) && typeof value === 'object' && value && value.__esModule) return value;
/******/ 		var ns = Object.create(null);
/******/ 		__webpack_require__.r(ns);
/******/ 		Object.defineProperty(ns, 'default', { enumerable: true, value: value });
/******/ 		if(mode & 2 && typeof value != 'string') for(var key in value) __webpack_require__.d(ns, key, function(key) { return value[key]; }.bind(null, key));
/******/ 		return ns;
/******/ 	};
/******/
/******/ 	// getDefaultExport function for compatibility with non-harmony modules
/******/ 	__webpack_require__.n = function(module) {
/******/ 		var getter = module && module.__esModule ?
/******/ 			function getDefault() { return module['default']; } :
/******/ 			function getModuleExports() { return module; };
/******/ 		__webpack_require__.d(getter, 'a', getter);
/******/ 		return getter;
/******/ 	};
/******/
/******/ 	// Object.prototype.hasOwnProperty.call
/******/ 	__webpack_require__.o = function(object, property) { return Object.prototype.hasOwnProperty.call(object, property); };
/******/
/******/ 	// __webpack_public_path__
/******/ 	__webpack_require__.p = "/static/";
/******/
/******/
/******/ 	// Load entry module and return exports
/******/ 	return __webpack_require__(__webpack_require__.s = "./client/main.js");
/******/ })
/************************************************************************/
/******/ ({

/***/ "./client/main.js":
/*!************************!*\
  !*** ./client/main.js ***!
  \************************/
/*! no static exports found */
/***/ (function(module, exports) {

eval("const DEFAULT_WAIT_MS = 1000\nconst MAX_MESSAGES = 1000\nconst KEY_ENTER = 'Enter'\n\nfunction postData (url = '', data = {}) {\n  return fetch(url, {\n    method: 'POST',\n    headers: {\n      'Content-Type': 'application/json; charset=utf-8'\n    },\n    body: JSON.stringify(data)\n  })\n}\n\nfunction getMessage () {\n  return document.getElementById('message')\n}\n\nfunction handleMessageKey (e) {\n  if (e.key === KEY_ENTER) {\n    e.preventDefault()\n    submitMessage()\n  }\n}\n\nfunction setupSubmit () {\n  const msgElem = getMessage()\n  msgElem.onkeyup = handleMessageKey\n  const sendElem = document.getElementById('send')\n  sendElem.onclick = submitMessage\n}\n\nasync function submitMessage () {\n  const elem = getMessage()\n  const message = elem.value\n  console.debug(`Submitting message \"${message}\"`)\n  elem.value = ''\n  await postData('/messages/', { body: message })\n}\n\nfunction appendMessages (messages = [], buff = []) {\n  const content = messages.map(m => m.body)\n  if (content.length > 0) {\n    console.debug(`Appending ${content.length} total messages.`)\n    const total = content.length + buff.length\n    const deficit = MAX_MESSAGES - total\n    if (deficit < 0) {\n      buff = buff.slice(0 - deficit).concat(content)\n    }\n  } else {\n    console.debug('No messages')\n  }\n  return buff\n}\n\nfunction getChat () {\n  return document.getElementById('chat')\n}\n\nfunction writeMessages (buff) {\n  const elem = getChat()\n  elem.value = buff.join('\\n')\n}\n\nasync function sleeper (ms = 1000) {\n  return new Promise((resolve, reject) => {\n    setTimeout(() => {\n      resolve()\n    }, ms)\n  })\n}\n\nasync function * loadMessages () {\n  let failure = 0\n  let ms = DEFAULT_WAIT_MS\n  while (true && failure < 10) {\n    const wait = sleeper(ms)\n    const resp = await fetch('/messages/')\n    yield resp\n    if (resp.status !== 200) {\n      failure++\n      ms *= 2\n    } else {\n      failure = 0\n      ms = DEFAULT_WAIT_MS\n    }\n    await wait\n  }\n}\n\n(async function main () {\n  setupSubmit()\n\n  let buff = []\n  for await (const resp of loadMessages()) {\n    console.debug(`Received ${resp.status} response.`)\n    buff = appendMessages(await resp.json(), buff)\n    writeMessages(buff)\n  }\n  console.debug('Giving up.')\n})()\n\n\n//# sourceURL=webpack:///./client/main.js?");

/***/ })

/******/ });