//import pjson from 'pjson'

import React from 'react'
import { render } from 'react-dom'
import styled, { keyframes } from 'styled-components'



type TitleProps = {
   title?: string
}

const Titolo = styled.h1`
  font-size: 1.5em;
  text-align: left;
  color: white;
`

const Wrapper = styled.section`
padding: 1em;
background: black;
`
const Flag = styled.div`
display: flex;
flex-direction: row;
`

const Green = styled.div`
width: 30vh;
background: green;
flex-direction: row;
`

const White = styled.div`
width: 30vh;
background: white;
flex-direction: row;
`

const Red = styled.div`
width: 30vh;
background: red;
flex-direction: row;
`


export default function Title({ title = "camarones" }) {
   return (
      <>
         <div>
            <Titolo>{title}</Titolo>
            <Flag>
               <div><Green>green</Green></div>
               <div><White>white</White></div>
               <div><Red>red</Red></div>
            </Flag>
         </div>
      </>
   )
}
