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
padding: 4em;
background: black;
`
const Flag = styled.section`
flex-direction: row;
`

const Green = styled.section`
background: green;
flex-direction: column;
`

const White = styled.section`
background: white;
flex-direction: column;
`

const Red = styled.section`
background: red;
flex-direction: column;
`


export default function Title({ title = "hatching" }) {
   return (
      <>
         <div>
            <Wrapper>
            <Titolo>{title}</Titolo>
            This apparat can help penguin's eggs to grow up!
            <Flag>
               <div><Green>green</Green></div>
               <div><White>white</White></div>
               <div><Red>red</Red></div>
            </Flag>
            </Wrapper>
         </div>
      </>
   )
}
