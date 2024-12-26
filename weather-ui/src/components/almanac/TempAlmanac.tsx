import BoxAlmanac from "./BoxAlmanac.tsx";
import LineChart from "./LineChart.tsx";
const TempAlmanac = () => {

    return (
       <BoxAlmanac>
           <>
               <LineChart apiUrl="http://10.10.1.83:5173/api/chart/tempf/day" />
           </>
       </BoxAlmanac>
    )
}
export default TempAlmanac
