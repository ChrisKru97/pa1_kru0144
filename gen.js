n = 13
arr=[]
for(let i = 0; i < n; i++) {
    arr.push([])
    for(let j = 0; j < n; j++) {
        const value = Math.round(Math.random()*100)
        if(j === i) arr[i].push(0);
        else if(i < j) {
            arr[i].push(value);
        } else {
            arr[i].push(arr[j][i])
        }
    }
}
let returnString = "{\n"
arr.forEach(i=> {
	returnString+="{"
	i.forEach(j=>returnString+=`${j},`)
	returnString+="},\n"
});
returnString+="}"

console.log(returnString)
