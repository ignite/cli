export const getters = {}

getters.jsonToHtmlTree = (jsonObj) => {
    const keys = []
    let treeHolder = ''

    for (let key in jsonObj) {
        if (typeof jsonObj[key] === 'object') {
          treeHolder += `<div class="wrapper">`
          treeHolder += `<span class="wrapper__key-item">${key}:</span>`
          treeHolder += `<div class="wrapper">${getters.jsonToHtmlTree(jsonObj[key])}</div>`
          treeHolder += `</div>`
        } else {
          treeHolder += `<span class="wrapper__item">${key}: ${jsonObj[key]}</span>`
        }
        keys.push(key)
    }
    return treeHolder  
}