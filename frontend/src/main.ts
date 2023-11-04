import {HttpClientModule} from '@angular/common/http';
import {importProvidersFrom} from '@angular/core';
import {ReactiveFormsModule} from '@angular/forms';
import {BrowserModule, bootstrapApplication} from '@angular/platform-browser';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {provideRouter} from '@angular/router';

import {AppComponent} from './app/app.component';
import {appRoutes} from './app/app.routes';

bootstrapApplication(AppComponent, {
    providers: [importProvidersFrom(BrowserModule, BrowserAnimationsModule, ReactiveFormsModule, HttpClientModule), provideRouter(appRoutes)]
}).catch(err => console.error(err));
