import {Component} from '@angular/core';
import {RouterOutlet} from '@angular/router';

import {NgParticlesModule} from 'ng-particles';
import {Engine, MoveDirection} from 'tsparticles-engine';
import {loadSnowPreset} from 'tsparticles-preset-snow';

@Component({
    standalone: true,
    selector: 'app-root',
    templateUrl: './app.component.html',
    imports: [RouterOutlet, NgParticlesModule],
    styleUrls: ['./app.component.scss']
})
export class AppComponent {
    id = 'tsparticles';

    async particlesInit(engine: Engine): Promise<void> {
        await loadSnowPreset(engine);
    }

    // These are all options copied from preset-snow
    particlesOptions = {
        background: {
            color: '',
            opacity: 1
        },
        particles: {
            number: {
                value: 200
            },
            move: {
                direction: MoveDirection.bottom,
                enable: true,
                random: false,
                straight: true,
                speed: 1
            },
            opacity: {
                value: {min: 0.1, max: 0.9}
            },
            size: {
                value: {min: 0.5, max: 5}
            },
            wobble: {
                distance: 1,
                enable: true,
                speed: {
                    min: 1,
                    max: 1
                }
            }
        }
    };
}
