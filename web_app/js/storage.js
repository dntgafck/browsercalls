/* More information about these options at jshint.com/docs/options */

/* exported Storage */
/* globals isChromeApp, chrome */

'use strict';

var Storage = function() {};

// Get a value from local browser storage. Calls callback with value.
// Handles variation in API between localStorage and Chrome app storage.
Storage.prototype.getStorage = function(key, callback) {
  if (isChromeApp()) {
    // Use chrome.storage.local.
    chrome.storage.local.get(key, function(values) {
      // Unwrap key/value pair.
      if (callback) {
        window.setTimeout(function() {
          callback(values[key]);
        }, 0);
      }
    });
  } else {
    // Use localStorage.
    var value = localStorage.getItem(key);
    if (callback) {
      window.setTimeout(function() {
        callback(value);
      }, 0);
    }
  }
};

// Set a value in local browser storage. Calls callback after completion.
// Handles variation in API between localStorage and Chrome app storage.
Storage.prototype.setStorage = function(key, value, callback) {
  if (isChromeApp()) {
    // Use chrome.storage.local.
    var data = {};
    data[key] = value;
    chrome.storage.local.set(data, callback);
  } else {
    // Use localStorage.
    localStorage.setItem(key, value);
    if (callback) {
      window.setTimeout(callback, 0);
    }
  }
};
