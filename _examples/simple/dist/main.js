"use strict";(()=>{function r(e){let o=0,t=i=>{o=i,e.innerHTML=`count is ${o}`};e.addEventListener("click",()=>t(o+1)),t(0)}document.querySelector("#app").innerHTML=`
  <div>
    <h1>TypeScript + Go</h1>
    <div class="card">
      <button id="counter" type="button"></button>
    </div>
  </div>
`;r(document.querySelector("#counter"));})();
//# sourceMappingURL=main.js.map
