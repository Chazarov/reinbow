/*
 * ATTENTION: The "eval" devtool has been used (maybe by default in mode: "development").
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
/******/ (() => { // webpackBootstrap
/******/ 	"use strict";
/******/ 	var __webpack_modules__ = ({

/***/ "./node_modules/css-loader/dist/cjs.js!./static/styles.css"
/*!*****************************************************************!*\
  !*** ./node_modules/css-loader/dist/cjs.js!./static/styles.css ***!
  \*****************************************************************/
(module, __webpack_exports__, __webpack_require__) {

eval("{__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"default\": () => (__WEBPACK_DEFAULT_EXPORT__)\n/* harmony export */ });\n/* harmony import */ var _node_modules_css_loader_dist_runtime_noSourceMaps_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ../node_modules/css-loader/dist/runtime/noSourceMaps.js */ \"./node_modules/css-loader/dist/runtime/noSourceMaps.js\");\n/* harmony import */ var _node_modules_css_loader_dist_runtime_noSourceMaps_js__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(_node_modules_css_loader_dist_runtime_noSourceMaps_js__WEBPACK_IMPORTED_MODULE_0__);\n/* harmony import */ var _node_modules_css_loader_dist_runtime_api_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ../node_modules/css-loader/dist/runtime/api.js */ \"./node_modules/css-loader/dist/runtime/api.js\");\n/* harmony import */ var _node_modules_css_loader_dist_runtime_api_js__WEBPACK_IMPORTED_MODULE_1___default = /*#__PURE__*/__webpack_require__.n(_node_modules_css_loader_dist_runtime_api_js__WEBPACK_IMPORTED_MODULE_1__);\n// Imports\n\n\nvar ___CSS_LOADER_EXPORT___ = _node_modules_css_loader_dist_runtime_api_js__WEBPACK_IMPORTED_MODULE_1___default()((_node_modules_css_loader_dist_runtime_noSourceMaps_js__WEBPACK_IMPORTED_MODULE_0___default()));\n// Module\n___CSS_LOADER_EXPORT___.push([module.id, `.container\n{\n    display: flex;\n    width: 70%;\n    height: 25%;\n    justify-content: center;\n}\n\n.headerArea\n{\n    width: 30%;\n    border-bottom: 4px solid transparent;\n    border-image: linear-gradient(to right, rgb(0, 132, 255), rgb(0, 204, 255)) 1;\n}\n\n.webhookTable\n{\n    margin-top: 20px;\n    width: 33%;\n}\n\n.webhookTable td:last-child {\n    width: 250px;\n}\n\n.containerWebhooks\n{\n    margin-bottom: 30px;\n}\n\n.text \n{\n    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;\n}\n\n.text-bold\n{\n    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;\n    font-weight: 600;\n}\n\n.table_main\n{\n    width:100%; \n    border-collapse:collapse; \n    text-align:center; \n    font-family:sans-serif; \n    font-size:14px;\n}\n\n.option\n{\n    margin-top: 25px;\n    margin-bottom: 15px;\n}\n\n/* Кастомный селект для таблицы */\n.custom-select {\n    position: relative;\n    width: 100%;\n    max-width: 200px;\n    display: inline-block;\n    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;\n    user-select: none;\n}\n\n.select-display {\n    border: 1px solid #ccc;\n    border-radius: 4px;\n    padding: 6px 10px;\n    background: white;\n    cursor: pointer;\n    display: flex;\n    justify-content: space-between;\n    align-items: center;\n    min-height: 34px;\n    font-size: 14px;\n}\n\n.select-placeholder {\n    color: #999;\n    font-size: 14px;\n    white-space: nowrap;\n    overflow: hidden;\n    text-overflow: ellipsis;\n    flex: 1;\n}\n\n.select-arrow {\n    color: #555;\n    font-size: 14px;\n    margin-left: 8px;\n    transition: transform 0.2s;\n}\n\n.select-options {\n    position: absolute;\n    top: 100%;\n    left: 0;\n    right: 0;\n    margin-top: 4px;\n    background: white;\n    border: 1px solid #ccc;\n    border-radius: 4px;\n    box-shadow: 0 4px 8px rgba(0,0,0,0.1);\n    max-height: 200px;\n    overflow-y: auto;\n    z-index: 10;\n    padding: 4px 0;\n    min-width: 180px;\n}\n\n.option-item {\n    display: flex;\n    align-items: center;\n    padding: 6px 12px;\n    cursor: pointer;\n    font-size: 14px;\n    transition: background 0.15s;\n}\n\n.option-item:hover {\n    background: #f0f0f0;\n}\n\n.option-item input[type=\"checkbox\"] {\n    margin-right: 8px;\n    cursor: pointer;\n    flex-shrink: 0;\n}\n\n.custom-select.open .select-arrow {\n    transform: rotate(180deg);\n}\n\n.url_input\n{\n    width: 80%;\n    height: 35px;\n    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;\n}\n\n.buttonArea\n{\n}\n\n.btnSaveSettings\n{\n    width: 100px;\n    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;\n}\n\n.modal_login\n{\n    width: 65%;\n    height: 3%;\n}\n\n.modal_password\n{\n    width: 65%;\n    height: 3%;\n}\n\n.krestik\n{\n    right: 0%;\n    position: absolute;\n    width: 7%;\n    color: white;\n    background-color: red;\n    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;\n    font-weight: bold;\n    font-size: small;\n    cursor: pointer;\n}\n\n.login_text\n{\n    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;\n    font-weight: 600;\n}\n\n.login-btn\n{\n    margin: 10px;\n    width: 25%;\n    height: 4%;\n}\n\n.black_overlay {\n    display: block;\n    position: fixed;\n    top: 0;\n    left: 0;\n    width: 100vw;\n    height: 100vh;\n    background: rgba(0, 0, 0, 0);\n    z-index: 1;      /* больше, чем у всех остальных элементов */\n    /*display: none;      /* изначально скрыт */\n    /* опционально – запрещаем клики сквозь слой */\n    pointer-events: none;\n    transition: all 0.5s ease-in-out; \n}\n\n.modal_container\n{\n    opacity: 0%;\n    border-radius: 10px;\n    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.5), 0 4px 6px -2px rgba(0, 0, 0, 0.5);\n    overflow: hidden;\n    position:fixed;\n    top: 50%;\n    left: 50%;\n    transform: translate(-50%, -50%);\n    display: block;\n    z-index: -1;\n    margin-top: 10px;\n    width: 20%;\n    background-color: lightgray;\n    transition: all 0.5s ease-in-out; \n}\n\n.authorize\n{\n    position: fixed;\n    top: 0%;\n    right: 0%;\n    transform: translate(-20%, 50%);\n}\n\n.userAuthorized\n{\n    position: fixed;\n    top: 0%;\n    right: 0%;\n    transform: translate(-40%, -10%);\n    display: none;\n}\n\n.btnLogout\n{\n    margin-top: -10px;\n}\n\n.errorNotification\n{\n    position: fixed;\n    bottom: 0%;\n    left: 0%;\n    transform: translate(12%, 100%);\n    border-radius: 5px;\n    overflow: hidden;\n    background-color: rgba(255, 0, 0, 0.75);\n    width: 15%;\n    height: 65px;\n    transition: transform 0.5s ease-in-out; \n    z-index: 2;\n}\n\n.typeOfNotification\n{\n    color: white;\n    background-color: rgba(0, 0, 0, 0.5);\n}\n\n.errorText\n{\n    margin-left: 10px;\n    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;\n}\n\n.errDesc\n{\n    display: flex;\n    margin-top: 7px;\n    margin-left: 10px;\n    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;\n    font-size: 12pt;\n    color: black;\n}`, \"\"]);\n// Exports\n/* harmony default export */ const __WEBPACK_DEFAULT_EXPORT__ = (___CSS_LOADER_EXPORT___);\n\n\n//# sourceURL=webpack://rbs-study-proj/./static/styles.css?./node_modules/css-loader/dist/cjs.js\n}");

/***/ },

/***/ "./node_modules/css-loader/dist/runtime/api.js"
/*!*****************************************************!*\
  !*** ./node_modules/css-loader/dist/runtime/api.js ***!
  \*****************************************************/
(module) {

eval("{\n\n/*\n  MIT License http://www.opensource.org/licenses/mit-license.php\n  Author Tobias Koppers @sokra\n*/\nmodule.exports = function (cssWithMappingToString) {\n  var list = [];\n\n  // return the list of modules as css string\n  list.toString = function toString() {\n    return this.map(function (item) {\n      var content = \"\";\n      var needLayer = typeof item[5] !== \"undefined\";\n      if (item[4]) {\n        content += \"@supports (\".concat(item[4], \") {\");\n      }\n      if (item[2]) {\n        content += \"@media \".concat(item[2], \" {\");\n      }\n      if (needLayer) {\n        content += \"@layer\".concat(item[5].length > 0 ? \" \".concat(item[5]) : \"\", \" {\");\n      }\n      content += cssWithMappingToString(item);\n      if (needLayer) {\n        content += \"}\";\n      }\n      if (item[2]) {\n        content += \"}\";\n      }\n      if (item[4]) {\n        content += \"}\";\n      }\n      return content;\n    }).join(\"\");\n  };\n\n  // import a list of modules into the list\n  list.i = function i(modules, media, dedupe, supports, layer) {\n    if (typeof modules === \"string\") {\n      modules = [[null, modules, undefined]];\n    }\n    var alreadyImportedModules = {};\n    if (dedupe) {\n      for (var k = 0; k < this.length; k++) {\n        var id = this[k][0];\n        if (id != null) {\n          alreadyImportedModules[id] = true;\n        }\n      }\n    }\n    for (var _k = 0; _k < modules.length; _k++) {\n      var item = [].concat(modules[_k]);\n      if (dedupe && alreadyImportedModules[item[0]]) {\n        continue;\n      }\n      if (typeof layer !== \"undefined\") {\n        if (typeof item[5] === \"undefined\") {\n          item[5] = layer;\n        } else {\n          item[1] = \"@layer\".concat(item[5].length > 0 ? \" \".concat(item[5]) : \"\", \" {\").concat(item[1], \"}\");\n          item[5] = layer;\n        }\n      }\n      if (media) {\n        if (!item[2]) {\n          item[2] = media;\n        } else {\n          item[1] = \"@media \".concat(item[2], \" {\").concat(item[1], \"}\");\n          item[2] = media;\n        }\n      }\n      if (supports) {\n        if (!item[4]) {\n          item[4] = \"\".concat(supports);\n        } else {\n          item[1] = \"@supports (\".concat(item[4], \") {\").concat(item[1], \"}\");\n          item[4] = supports;\n        }\n      }\n      list.push(item);\n    }\n  };\n  return list;\n};\n\n//# sourceURL=webpack://rbs-study-proj/./node_modules/css-loader/dist/runtime/api.js?\n}");

/***/ },

/***/ "./node_modules/css-loader/dist/runtime/noSourceMaps.js"
/*!**************************************************************!*\
  !*** ./node_modules/css-loader/dist/runtime/noSourceMaps.js ***!
  \**************************************************************/
(module) {

eval("{\n\nmodule.exports = function (i) {\n  return i[1];\n};\n\n//# sourceURL=webpack://rbs-study-proj/./node_modules/css-loader/dist/runtime/noSourceMaps.js?\n}");

/***/ },

/***/ "./static/styles.css"
/*!***************************!*\
  !*** ./static/styles.css ***!
  \***************************/
(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

eval("{__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"default\": () => (__WEBPACK_DEFAULT_EXPORT__)\n/* harmony export */ });\n/* harmony import */ var _node_modules_style_loader_dist_runtime_injectStylesIntoStyleTag_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! !../node_modules/style-loader/dist/runtime/injectStylesIntoStyleTag.js */ \"./node_modules/style-loader/dist/runtime/injectStylesIntoStyleTag.js\");\n/* harmony import */ var _node_modules_style_loader_dist_runtime_injectStylesIntoStyleTag_js__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(_node_modules_style_loader_dist_runtime_injectStylesIntoStyleTag_js__WEBPACK_IMPORTED_MODULE_0__);\n/* harmony import */ var _node_modules_style_loader_dist_runtime_styleDomAPI_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! !../node_modules/style-loader/dist/runtime/styleDomAPI.js */ \"./node_modules/style-loader/dist/runtime/styleDomAPI.js\");\n/* harmony import */ var _node_modules_style_loader_dist_runtime_styleDomAPI_js__WEBPACK_IMPORTED_MODULE_1___default = /*#__PURE__*/__webpack_require__.n(_node_modules_style_loader_dist_runtime_styleDomAPI_js__WEBPACK_IMPORTED_MODULE_1__);\n/* harmony import */ var _node_modules_style_loader_dist_runtime_insertBySelector_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! !../node_modules/style-loader/dist/runtime/insertBySelector.js */ \"./node_modules/style-loader/dist/runtime/insertBySelector.js\");\n/* harmony import */ var _node_modules_style_loader_dist_runtime_insertBySelector_js__WEBPACK_IMPORTED_MODULE_2___default = /*#__PURE__*/__webpack_require__.n(_node_modules_style_loader_dist_runtime_insertBySelector_js__WEBPACK_IMPORTED_MODULE_2__);\n/* harmony import */ var _node_modules_style_loader_dist_runtime_setAttributesWithoutAttributes_js__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! !../node_modules/style-loader/dist/runtime/setAttributesWithoutAttributes.js */ \"./node_modules/style-loader/dist/runtime/setAttributesWithoutAttributes.js\");\n/* harmony import */ var _node_modules_style_loader_dist_runtime_setAttributesWithoutAttributes_js__WEBPACK_IMPORTED_MODULE_3___default = /*#__PURE__*/__webpack_require__.n(_node_modules_style_loader_dist_runtime_setAttributesWithoutAttributes_js__WEBPACK_IMPORTED_MODULE_3__);\n/* harmony import */ var _node_modules_style_loader_dist_runtime_insertStyleElement_js__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! !../node_modules/style-loader/dist/runtime/insertStyleElement.js */ \"./node_modules/style-loader/dist/runtime/insertStyleElement.js\");\n/* harmony import */ var _node_modules_style_loader_dist_runtime_insertStyleElement_js__WEBPACK_IMPORTED_MODULE_4___default = /*#__PURE__*/__webpack_require__.n(_node_modules_style_loader_dist_runtime_insertStyleElement_js__WEBPACK_IMPORTED_MODULE_4__);\n/* harmony import */ var _node_modules_style_loader_dist_runtime_styleTagTransform_js__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(/*! !../node_modules/style-loader/dist/runtime/styleTagTransform.js */ \"./node_modules/style-loader/dist/runtime/styleTagTransform.js\");\n/* harmony import */ var _node_modules_style_loader_dist_runtime_styleTagTransform_js__WEBPACK_IMPORTED_MODULE_5___default = /*#__PURE__*/__webpack_require__.n(_node_modules_style_loader_dist_runtime_styleTagTransform_js__WEBPACK_IMPORTED_MODULE_5__);\n/* harmony import */ var _node_modules_css_loader_dist_cjs_js_styles_css__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(/*! !!../node_modules/css-loader/dist/cjs.js!./styles.css */ \"./node_modules/css-loader/dist/cjs.js!./static/styles.css\");\n\n      \n      \n      \n      \n      \n      \n      \n      \n      \n\nvar options = {};\n\noptions.styleTagTransform = (_node_modules_style_loader_dist_runtime_styleTagTransform_js__WEBPACK_IMPORTED_MODULE_5___default());\noptions.setAttributes = (_node_modules_style_loader_dist_runtime_setAttributesWithoutAttributes_js__WEBPACK_IMPORTED_MODULE_3___default());\noptions.insert = _node_modules_style_loader_dist_runtime_insertBySelector_js__WEBPACK_IMPORTED_MODULE_2___default().bind(null, \"head\");\noptions.domAPI = (_node_modules_style_loader_dist_runtime_styleDomAPI_js__WEBPACK_IMPORTED_MODULE_1___default());\noptions.insertStyleElement = (_node_modules_style_loader_dist_runtime_insertStyleElement_js__WEBPACK_IMPORTED_MODULE_4___default());\n\nvar update = _node_modules_style_loader_dist_runtime_injectStylesIntoStyleTag_js__WEBPACK_IMPORTED_MODULE_0___default()(_node_modules_css_loader_dist_cjs_js_styles_css__WEBPACK_IMPORTED_MODULE_6__[\"default\"], options);\n\n\n\n\n       /* harmony default export */ const __WEBPACK_DEFAULT_EXPORT__ = (_node_modules_css_loader_dist_cjs_js_styles_css__WEBPACK_IMPORTED_MODULE_6__[\"default\"] && _node_modules_css_loader_dist_cjs_js_styles_css__WEBPACK_IMPORTED_MODULE_6__[\"default\"].locals ? _node_modules_css_loader_dist_cjs_js_styles_css__WEBPACK_IMPORTED_MODULE_6__[\"default\"].locals : undefined);\n\n\n//# sourceURL=webpack://rbs-study-proj/./static/styles.css?\n}");

/***/ },

/***/ "./node_modules/style-loader/dist/runtime/injectStylesIntoStyleTag.js"
/*!****************************************************************************!*\
  !*** ./node_modules/style-loader/dist/runtime/injectStylesIntoStyleTag.js ***!
  \****************************************************************************/
(module) {

eval("{\n\nvar stylesInDOM = [];\nfunction getIndexByIdentifier(identifier) {\n  var result = -1;\n  for (var i = 0; i < stylesInDOM.length; i++) {\n    if (stylesInDOM[i].identifier === identifier) {\n      result = i;\n      break;\n    }\n  }\n  return result;\n}\nfunction modulesToDom(list, options) {\n  var idCountMap = {};\n  var identifiers = [];\n  for (var i = 0; i < list.length; i++) {\n    var item = list[i];\n    var id = options.base ? item[0] + options.base : item[0];\n    var count = idCountMap[id] || 0;\n    var identifier = \"\".concat(id, \" \").concat(count);\n    idCountMap[id] = count + 1;\n    var indexByIdentifier = getIndexByIdentifier(identifier);\n    var obj = {\n      css: item[1],\n      media: item[2],\n      sourceMap: item[3],\n      supports: item[4],\n      layer: item[5]\n    };\n    if (indexByIdentifier !== -1) {\n      stylesInDOM[indexByIdentifier].references++;\n      stylesInDOM[indexByIdentifier].updater(obj);\n    } else {\n      var updater = addElementStyle(obj, options);\n      options.byIndex = i;\n      stylesInDOM.splice(i, 0, {\n        identifier: identifier,\n        updater: updater,\n        references: 1\n      });\n    }\n    identifiers.push(identifier);\n  }\n  return identifiers;\n}\nfunction addElementStyle(obj, options) {\n  var api = options.domAPI(options);\n  api.update(obj);\n  var updater = function updater(newObj) {\n    if (newObj) {\n      if (newObj.css === obj.css && newObj.media === obj.media && newObj.sourceMap === obj.sourceMap && newObj.supports === obj.supports && newObj.layer === obj.layer) {\n        return;\n      }\n      api.update(obj = newObj);\n    } else {\n      api.remove();\n    }\n  };\n  return updater;\n}\nmodule.exports = function (list, options) {\n  options = options || {};\n  list = list || [];\n  var lastIdentifiers = modulesToDom(list, options);\n  return function update(newList) {\n    newList = newList || [];\n    for (var i = 0; i < lastIdentifiers.length; i++) {\n      var identifier = lastIdentifiers[i];\n      var index = getIndexByIdentifier(identifier);\n      stylesInDOM[index].references--;\n    }\n    var newLastIdentifiers = modulesToDom(newList, options);\n    for (var _i = 0; _i < lastIdentifiers.length; _i++) {\n      var _identifier = lastIdentifiers[_i];\n      var _index = getIndexByIdentifier(_identifier);\n      if (stylesInDOM[_index].references === 0) {\n        stylesInDOM[_index].updater();\n        stylesInDOM.splice(_index, 1);\n      }\n    }\n    lastIdentifiers = newLastIdentifiers;\n  };\n};\n\n//# sourceURL=webpack://rbs-study-proj/./node_modules/style-loader/dist/runtime/injectStylesIntoStyleTag.js?\n}");

/***/ },

/***/ "./node_modules/style-loader/dist/runtime/insertBySelector.js"
/*!********************************************************************!*\
  !*** ./node_modules/style-loader/dist/runtime/insertBySelector.js ***!
  \********************************************************************/
(module) {

eval("{\n\nvar memo = {};\n\n/* istanbul ignore next  */\nfunction getTarget(target) {\n  if (typeof memo[target] === \"undefined\") {\n    var styleTarget = document.querySelector(target);\n\n    // Special case to return head of iframe instead of iframe itself\n    if (window.HTMLIFrameElement && styleTarget instanceof window.HTMLIFrameElement) {\n      try {\n        // This will throw an exception if access to iframe is blocked\n        // due to cross-origin restrictions\n        styleTarget = styleTarget.contentDocument.head;\n      } catch (e) {\n        // istanbul ignore next\n        styleTarget = null;\n      }\n    }\n    memo[target] = styleTarget;\n  }\n  return memo[target];\n}\n\n/* istanbul ignore next  */\nfunction insertBySelector(insert, style) {\n  var target = getTarget(insert);\n  if (!target) {\n    throw new Error(\"Couldn't find a style target. This probably means that the value for the 'insert' parameter is invalid.\");\n  }\n  target.appendChild(style);\n}\nmodule.exports = insertBySelector;\n\n//# sourceURL=webpack://rbs-study-proj/./node_modules/style-loader/dist/runtime/insertBySelector.js?\n}");

/***/ },

/***/ "./node_modules/style-loader/dist/runtime/insertStyleElement.js"
/*!**********************************************************************!*\
  !*** ./node_modules/style-loader/dist/runtime/insertStyleElement.js ***!
  \**********************************************************************/
(module) {

eval("{\n\n/* istanbul ignore next  */\nfunction insertStyleElement(options) {\n  var element = document.createElement(\"style\");\n  options.setAttributes(element, options.attributes);\n  options.insert(element, options.options);\n  return element;\n}\nmodule.exports = insertStyleElement;\n\n//# sourceURL=webpack://rbs-study-proj/./node_modules/style-loader/dist/runtime/insertStyleElement.js?\n}");

/***/ },

/***/ "./node_modules/style-loader/dist/runtime/setAttributesWithoutAttributes.js"
/*!**********************************************************************************!*\
  !*** ./node_modules/style-loader/dist/runtime/setAttributesWithoutAttributes.js ***!
  \**********************************************************************************/
(module, __unused_webpack_exports, __webpack_require__) {

eval("{\n\n/* istanbul ignore next  */\nfunction setAttributesWithoutAttributes(styleElement) {\n  var nonce =  true ? __webpack_require__.nc : 0;\n  if (nonce) {\n    styleElement.setAttribute(\"nonce\", nonce);\n  }\n}\nmodule.exports = setAttributesWithoutAttributes;\n\n//# sourceURL=webpack://rbs-study-proj/./node_modules/style-loader/dist/runtime/setAttributesWithoutAttributes.js?\n}");

/***/ },

/***/ "./node_modules/style-loader/dist/runtime/styleDomAPI.js"
/*!***************************************************************!*\
  !*** ./node_modules/style-loader/dist/runtime/styleDomAPI.js ***!
  \***************************************************************/
(module) {

eval("{\n\n/* istanbul ignore next  */\nfunction apply(styleElement, options, obj) {\n  var css = \"\";\n  if (obj.supports) {\n    css += \"@supports (\".concat(obj.supports, \") {\");\n  }\n  if (obj.media) {\n    css += \"@media \".concat(obj.media, \" {\");\n  }\n  var needLayer = typeof obj.layer !== \"undefined\";\n  if (needLayer) {\n    css += \"@layer\".concat(obj.layer.length > 0 ? \" \".concat(obj.layer) : \"\", \" {\");\n  }\n  css += obj.css;\n  if (needLayer) {\n    css += \"}\";\n  }\n  if (obj.media) {\n    css += \"}\";\n  }\n  if (obj.supports) {\n    css += \"}\";\n  }\n  var sourceMap = obj.sourceMap;\n  if (sourceMap && typeof btoa !== \"undefined\") {\n    css += \"\\n/*# sourceMappingURL=data:application/json;base64,\".concat(btoa(unescape(encodeURIComponent(JSON.stringify(sourceMap)))), \" */\");\n  }\n\n  // For old IE\n  /* istanbul ignore if  */\n  options.styleTagTransform(css, styleElement, options.options);\n}\nfunction removeStyleElement(styleElement) {\n  // istanbul ignore if\n  if (styleElement.parentNode === null) {\n    return false;\n  }\n  styleElement.parentNode.removeChild(styleElement);\n}\n\n/* istanbul ignore next  */\nfunction domAPI(options) {\n  if (typeof document === \"undefined\") {\n    return {\n      update: function update() {},\n      remove: function remove() {}\n    };\n  }\n  var styleElement = options.insertStyleElement(options);\n  return {\n    update: function update(obj) {\n      apply(styleElement, options, obj);\n    },\n    remove: function remove() {\n      removeStyleElement(styleElement);\n    }\n  };\n}\nmodule.exports = domAPI;\n\n//# sourceURL=webpack://rbs-study-proj/./node_modules/style-loader/dist/runtime/styleDomAPI.js?\n}");

/***/ },

/***/ "./node_modules/style-loader/dist/runtime/styleTagTransform.js"
/*!*********************************************************************!*\
  !*** ./node_modules/style-loader/dist/runtime/styleTagTransform.js ***!
  \*********************************************************************/
(module) {

eval("{\n\n/* istanbul ignore next  */\nfunction styleTagTransform(css, styleElement) {\n  if (styleElement.styleSheet) {\n    styleElement.styleSheet.cssText = css;\n  } else {\n    while (styleElement.firstChild) {\n      styleElement.removeChild(styleElement.firstChild);\n    }\n    styleElement.appendChild(document.createTextNode(css));\n  }\n}\nmodule.exports = styleTagTransform;\n\n//# sourceURL=webpack://rbs-study-proj/./node_modules/style-loader/dist/runtime/styleTagTransform.js?\n}");

/***/ },

/***/ "./static/index.js"
/*!*************************!*\
  !*** ./static/index.js ***!
  \*************************/
(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

eval("{__webpack_require__.r(__webpack_exports__);\n/* harmony import */ var _styles_css__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./styles.css */ \"./static/styles.css\");\n  // ← добавляем эту строку\nconst sleep = ms => new Promise(r => setTimeout(r, ms));\n\n(function initNotifications() {\n    const HEIGHT = 65;             \n    const GAP = 20;                \n    const LIFETIME = 3000;         \n    const CONTAINER_WIDTH = '22%';\n    const LEFT_OFFSET = '2%';\n    const BOTTOM_OFFSET = 90;       \n\n    const template = document.getElementById('notification');\n    if (template) {\n        template.style.display = 'none';\n    }\n\n    const notifications = []; \n\n    function render() {\n        notifications.forEach((item, index) => {\n            const bottomPos = BOTTOM_OFFSET + index * (HEIGHT + GAP);\n            item.element.style.bottom = bottomPos + 'px';\n            item.element.style.opacity = '1';\n        });\n    }\n\n    function removeOldest() {\n        if (notifications.length === 0) return;\n\n        const oldest = notifications.pop();\n        const el = oldest.element;\n\n        el.style.transition = 'bottom 0.5s ease, opacity 0.5s ease';\n        el.style.bottom = (window.innerHeight + HEIGHT) + 'px';\n        el.style.opacity = '0';\n\n        const onFinish = () => {\n            el.remove();\n            el.removeEventListener('transitionend', onFinish);\n        };\n        el.addEventListener('transitionend', onFinish);\n\n        render(); \n    }\n\n    window.showNotification = function(type, title, description) {\n        if (!template) return;\n\n        const clone = template.cloneNode(true);\n        clone.id = '';\n        clone.style.display = 'block';\n        clone.style.position = 'fixed';\n        clone.style.left = LEFT_OFFSET;\n        clone.style.width = CONTAINER_WIDTH;\n        clone.style.height = HEIGHT + 'px';\n        clone.style.margin = '0';\n        clone.style.borderRadius = '5px';\n        clone.style.overflow = 'hidden';\n        clone.style.zIndex = '2';\n        clone.style.pointerEvents = 'none';\n\n        clone.style.backgroundColor = type === 'success'\n            ? 'rgba(0, 255, 0, 0.75)'\n            : 'rgba(255, 0, 0, 0.75)';\n\n        const titleSpan = clone.querySelector('.errorText');\n        const descSpan = clone.querySelector('.errDesc');\n        if (titleSpan) titleSpan.textContent = title;\n        if (descSpan) descSpan.textContent = description;\n\n        clone.style.transition = 'none';\n        clone.style.bottom = -(HEIGHT + 100) + 'px';\n        clone.style.opacity = '0';\n\n        notifications.unshift({ element: clone });\n\n        document.body.appendChild(clone);\n\n        clone.offsetHeight;\n\n        clone.style.transition = 'bottom 0.5s ease, opacity 0.5s ease';\n        render();\n\n        setTimeout(() => {\n            removeOldest();\n        }, LIFETIME);\n    };\n})();\n\nfunction getSelectedWebhookTypes() {\n    const select = document.getElementById('webhook_select_1');\n    // selectedOptions — коллекция выбранных <option>\n    return Array.from(select.selectedOptions).map(opt => opt.value);\n}\n\nasync function CompleteSetup() {\n    if (!isLogged) {\n        console.log(\"Необходимо войти в аккаунт!\");\n        showNotification('error', 'Ошибка', 'Необходимо войти в аккаунт!');\n        return;\n    }\n\n    let payload = {\n        \"data\": [\n            {\n                \"notify_type\": \"user_register\",\n                \"want_email\": document.getElementById('email_reg').checked,\n                \"want_telegram\": document.getElementById('tg_reg').checked,\n                \"want_webhook\": document.getElementById('webhook_reg').checked\n            },\n            {\n                \"notify_type\": \"user_login\",\n                \"want_email\": document.getElementById('email_login').checked,\n                \"want_telegram\": document.getElementById('tg_login').checked,\n                \"want_webhook\": document.getElementById('webhook_login').checked\n            },\n            {\n                \"notify_type\": \"admin_newImg\",\n                \"want_email\": document.getElementById('email_admin_newImg').checked,\n                \"want_telegram\": document.getElementById('tg_admin_newImg').checked,\n                \"want_webhook\": document.getElementById('webhook_admin_newImg').checked\n            }\n        ],\n        \"webhookData\": {\n            \"url\": document.getElementById('urlInput').value,\n            \"notificationTypes\": getSelectedWebhookTypes()\n        }\n    };\n\n    try {\n        const response = await fetch(\"http://10.64.10.237:8080/notify_types\", {\n            method: 'POST',\n            headers: { 'Content-Type': 'application/json' },\n            body: JSON.stringify(payload)\n        });\n\n        if (!response.ok) {\n            showNotification('error', 'Ошибка', 'При обращении к БД произошла ошибка')\n            throw new Error(`HTTP error! Status: ${response.status}`);\n        }\n        // Можно добавить успешное уведомление, если нужно\n        // showNotification('success', 'Успех', 'Настройки сохранены');\n    } catch (error) {\n        console.error('Error:', error);\n        showNotification('error', 'Ошибка', 'Не удалось сохранить настройки');\n    }\n}\n\nasync function GetSettings() {\n    try {\n        const response = await fetch(\"http://10.64.10.237:8080/get_notify_settings\");\n        const result = await response.json();\n        const data = result[\"data\"];\n        data.forEach((item) => {\n            const notify_type = item.notify_type;\n            const want_email = item.want_email;\n            const want_telegram = item.want_telegram;\n            const want_webhook = item.want_webhook\n\n            switch (notify_type) {\n                case \"user_register\":\n                    document.getElementById('email_reg').checked = want_email;\n                    document.getElementById('tg_reg').checked = want_telegram;\n                    document.getElementById('webhook_reg').checked = want_webhook;\n                    break;\n                case \"user_login\":\n                    document.getElementById('email_login').checked = want_email;\n                    document.getElementById('tg_login').checked = want_telegram;\n                    document.getElementById('webhook_login').checked = want_webhook;\n                    break;\n                case \"admin_newImg\":\n                    document.getElementById('email_admin_newImg').checked = want_email;\n                    document.getElementById('tg_admin_newImg').checked = want_telegram;\n                    document.getElementById('webhook_admin_newImg').checked = want_webhook;\n                    break;\n            }\n        });\n    } catch (error) {\n        console.error('Error:', error);\n    }\n}\n\nfunction closeModal() {\n    let modal_window = document.getElementById('modalWindow');\n    let black_background = document.getElementById('blackBackground');\n    modal_window.style.opacity = \"0%\";\n    modal_window.style.transform = \"translate(-50%, -50%)\";\n    black_background.style.pointerEvents = \"none\";\n    black_background.style.background = \"rgba(0, 0, 0, 0)\";\n    black_background.style.backdropFilter = \"none\";\n    modal_window.style.zIndex = \"-1\";\n}\n\nfunction openModal() {\n    let modal_window = document.getElementById('modalWindow');\n    let black_background = document.getElementById('blackBackground');\n    modal_window.style.zIndex = \"2\";\n    modal_window.style.opacity = \"100%\";\n    modal_window.style.transform = \"translate(-50%, -75%)\";\n    black_background.style.pointerEvents = \"auto\";\n    black_background.style.background = \"rgba(0, 0, 0, 0.5)\";\n    black_background.style.backdropFilter = \"blur(7px)\";\n}\n\nasync function loginProcedure() {\n    let login_value = document.getElementById('login_field').value;\n    let password_value = document.getElementById('password_field').value;\n\n    if (login_value == \"admin\" && password_value == \"12345\") {\n        console.log(\"Вход успешен!\");\n        isLogged = true;\n        document.getElementById('autorize').style.display = \"none\";\n        document.getElementById('isAuthorized').style.display = \"block\";\n        closeModal();\n        showNotification('success', 'Успех!', 'Успешный вход');\n    } else {\n        console.log(\"Неправильные данные!\");\n        showNotification('error', 'Ошибка', 'Неверно введен логин или пароль');\n    }\n}\n\nasync function logoutProcedure() {\n    document.getElementById('login_field').value = \"\";\n    document.getElementById('password_field').value = \"\";\n\n    document.getElementById('autorize').style.display = \"block\";\n    document.getElementById('isAuthorized').style.display = \"none\";\n\n    console.log(\"Выход успешен!\");\n    closeModal();\n    isLogged = false;\n    showNotification('success', 'Успех', 'Выход из аккаунта успешен!');\n}\n\nGetSettings();\nlet isLogged = false;\n\n(function initCustomSelect() {\n    const container = document.getElementById('webhookSelectContainer');\n    if (!container) return;\n\n    const display = container.querySelector('.select-display');\n    const optionsPanel = container.querySelector('.select-options');\n    const placeholder = display.querySelector('.select-placeholder');\n    const checkboxes = optionsPanel.querySelectorAll('input[type=\"checkbox\"]');\n    const hiddenSelect = document.getElementById('webhook_select_1');\n    const optionItems = optionsPanel.querySelectorAll('.option-item');\n\n    function updateDisplay() {\n        const checked = [];\n        const labels = [];\n        checkboxes.forEach(cb => {\n            if (cb.checked) {\n                const label = cb.closest('.option-item');\n                if (label) {\n                    labels.push(label.textContent.trim());\n                }\n                checked.push(cb.value);\n            }\n        });\n\n        if (checked.length === 0) {\n            placeholder.textContent = 'Ничего не выбрано';\n            placeholder.style.color = '#999';\n        } else {\n            placeholder.textContent = labels.join(', ');\n            placeholder.style.color = '#333';\n        }\n\n        // Синхронизация со скрытым select\n        if (hiddenSelect) {\n            Array.from(hiddenSelect.options).forEach(opt => {\n                opt.selected = checked.includes(opt.value);\n            });\n            hiddenSelect.dispatchEvent(new Event('change', { bubbles: true }));\n        }\n    }\n\n    // Открыть/закрыть список\n    display.addEventListener('click', function(e) {\n        e.stopPropagation();\n        const isOpen = container.classList.toggle('open');\n        optionsPanel.style.display = isOpen ? 'block' : 'none';\n    });\n\n    // Обработка клика по строке (включая текст и чекбокс)\n    optionItems.forEach(item => {\n        item.addEventListener('click', function(e) {\n            // Предотвращаем стандартное переключение чекбокса через label\n            e.preventDefault();\n\n            const checkbox = this.querySelector('input[type=\"checkbox\"]');\n            if (checkbox) {\n                // Переключаем состояние чекбокса\n                checkbox.checked = !checkbox.checked;\n                // Обновляем отображение\n                updateDisplay();\n            }\n        });\n    });\n\n    // Закрывать список при клике вне компонента\n    document.addEventListener('click', function(e) {\n        if (!container.contains(e.target)) {\n            container.classList.remove('open');\n            optionsPanel.style.display = 'none';\n        }\n    });\n\n    // Инициализация: ничего не выбрано\n    checkboxes.forEach(cb => cb.checked = false);\n    updateDisplay();\n})();\n\n//# sourceURL=webpack://rbs-study-proj/./static/index.js?\n}");

/***/ }

/******/ 	});
/************************************************************************/
/******/ 	// The module cache
/******/ 	const __webpack_module_cache__ = {};
/******/ 	
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/ 		// Check if module is in cache
/******/ 		const cachedModule = __webpack_module_cache__[moduleId];
/******/ 		if (cachedModule !== undefined) {
/******/ 			return cachedModule.exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		const module = __webpack_module_cache__[moduleId] = {
/******/ 			id: moduleId,
/******/ 			// no module.loaded needed
/******/ 			exports: {}
/******/ 		};
/******/ 	
/******/ 		// Execute the module function
/******/ 		if (!(moduleId in __webpack_modules__)) {
/******/ 			delete __webpack_module_cache__[moduleId];
/******/ 			const e = new Error("Cannot find module '" + moduleId + "'");
/******/ 			e.code = 'MODULE_NOT_FOUND';
/******/ 			throw e;
/******/ 		}
/******/ 		__webpack_modules__[moduleId](module, module.exports, __webpack_require__);
/******/ 	
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/ 	
/************************************************************************/
/******/ 	/* webpack/runtime/compat get default export */
/******/ 	(() => {
/******/ 		// getDefaultExport function for compatibility with non-harmony modules
/******/ 		__webpack_require__.n = (module) => {
/******/ 			const getter = module && module.__esModule ?
/******/ 				() => (module['default']) :
/******/ 				() => (module);
/******/ 			__webpack_require__.d(getter, { a: getter });
/******/ 			return getter;
/******/ 		};
/******/ 	})();
/******/ 	
/******/ 	/* webpack/runtime/define property getters */
/******/ 	(() => {
/******/ 		// define getter/value functions for harmony exports
/******/ 		__webpack_require__.d = (exports, definition) => {
/******/ 			if(Array.isArray(definition)) {
/******/ 				var i = 0;
/******/ 				while(i < definition.length) {
/******/ 					var key = definition[i++];
/******/ 					var binding = definition[i++];
/******/ 					if(!__webpack_require__.o(exports, key)) {
/******/ 						if(binding === 0) {
/******/ 							Object.defineProperty(exports, key, { enumerable: true, value: definition[i++] });
/******/ 						} else {
/******/ 							Object.defineProperty(exports, key, { enumerable: true, get: binding });
/******/ 						}
/******/ 					} else if(binding === 0) { i++; }
/******/ 				}
/******/ 			} else {
/******/ 				for(var key in definition) {
/******/ 					if(__webpack_require__.o(definition, key) && !__webpack_require__.o(exports, key)) {
/******/ 						Object.defineProperty(exports, key, { enumerable: true, get: definition[key] });
/******/ 					}
/******/ 				}
/******/ 			}
/******/ 		};
/******/ 	})();
/******/ 	
/******/ 	/* webpack/runtime/hasOwnProperty shorthand */
/******/ 	(() => {
/******/ 		__webpack_require__.o = (obj, prop) => (Object.prototype.hasOwnProperty.call(obj, prop))
/******/ 	})();
/******/ 	
/******/ 	/* webpack/runtime/make namespace object */
/******/ 	(() => {
/******/ 		// define __esModule on exports
/******/ 		__webpack_require__.r = (exports) => {
/******/ 			if(Symbol.toStringTag) {
/******/ 				Object.defineProperty(exports, Symbol.toStringTag, { value: 'Module' });
/******/ 			}
/******/ 			Object.defineProperty(exports, '__esModule', { value: true });
/******/ 		};
/******/ 	})();
/******/ 	
/******/ 	/* webpack/runtime/nonce */
/******/ 	(() => {
/******/ 		__webpack_require__.nc = undefined;
/******/ 	})();
/******/ 	
/************************************************************************/
/******/ 	
/******/ 	// startup
/******/ 	// Load entry module and return exports
/******/ 	// This entry module can't be inlined because the eval devtool is used.
/******/ 	let __webpack_exports__ = __webpack_require__("./static/index.js");
/******/ 	
/******/ })()
;