/**
 * Welcome
 */
 import React, { useState } from 'react'
 import {render} from 'react-dom'
 
 import pjson from 'pjson'
 import Title from './elements/title'
 import Steps from './elements/steps'
 
 // import yaml from 'js-yaml'
 // import fs from 'fs'
 import { ISettings, IBranding } from './interfaces'
 import { Container } from './styles'
 
 type WelcomeProps = {
   language?: string,
 }
 
 
 export default function Welcome({ language = '' }: WelcomeProps) {
   let productName = 'unknown'
   let version = 'x.x.x'
   let configRoot = '/etc/penguins-eggs.d/krill/'
   // if (fs.existsSync('/etc/calamares/settings.conf')) {
  // configRoot = '/etc/calamares/'
    //  }
 
   // const settings = yaml.load(fs.readFileSync(configRoot + 'settings.conf', 'utf-8')) as unknown as ISettings
   // const branding = settings.branding
   // const calamares = yaml.load(fs.readFileSync(configRoot + 'branding/' + branding + '/branding.desc', 'utf-8')) as unknown as IBranding
   // productName = calamares.strings.productName
   // version = calamares.strings.version
 
   return (
     <>
      <Title />
      <Steps step={2} />
     </>
   )
 }
 