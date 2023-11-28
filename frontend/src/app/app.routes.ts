import {Routes} from '@angular/router';

import {HeaderComponent} from './core/components/layout/header/header.component';
import {authGuard} from './core/guards/auth.guard';
import {AboutPageComponent} from './pages/about/about.page.component';
import {BingoBoardsPageComponent} from './pages/bingo-boards/bingo-boards.page.component';
import {BingoPageComponent} from './pages/bingo/bingo.page.component';
import {LoginPageComponent} from './pages/login/login.page.component';
import {PrivateBingoBoardPageComponent} from './pages/private-bingo-board/private-bingo-board.page.component';

export const appRoutes: Routes = [
    {
        path: '',
        component: HeaderComponent,
        children: [
            {
                path: '',
                component: BingoPageComponent
            },
            {
                path: 'about',
                component: AboutPageComponent
            },
            {
                path: 'private-bingo-boards',
                component: BingoBoardsPageComponent,
                canActivate: [authGuard]
            },
            {
                path: 'private-bingo-boards/:uuid',
                component: PrivateBingoBoardPageComponent,
                canActivate: [authGuard]
            },
            {
                path: 'login',
                component: LoginPageComponent
            }
        ]
    }
];
