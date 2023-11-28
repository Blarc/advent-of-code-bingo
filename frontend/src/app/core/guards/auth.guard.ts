import {inject} from '@angular/core';
import {CanActivateFn, Router, UrlTree} from '@angular/router';

import {Observable} from 'rxjs';

import {AuthService} from '../services/auth.service';

export const authGuard: CanActivateFn = (): Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree => {
    const router = inject(Router);
    const authenticationService = inject(AuthService);

    if (authenticationService.isTokenValid()) {
        return true;
    }
    return router.parseUrl('/');
};
